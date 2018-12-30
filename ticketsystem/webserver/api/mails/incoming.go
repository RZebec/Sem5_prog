package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"encoding/json"
	"github.com/pkg/errors"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
	A handler for incoming mails.
*/
type IncomingMailHandler struct {
	Logger      logging.Logger
	MailContext mail.MailContext
	TicketContext ticket.TicketContext
	UserContext user.UserContext
}

/*
	Handling the incoming mails.
*/
func (h *IncomingMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data []mail.Mail
	err := decoder.Decode(&data)
	if err != nil {
		h.Logger.LogError("IncomingMailHandler", err)
		w.WriteHeader(500)
	} else {
		err = h.handleIncomingMails(data)
		if err != nil {
			h.Logger.LogError("IncomingMailHandler", err)
			w.WriteHeader(500)
			return
		}
		h.Logger.LogInfo("IncomingMailHandler", "Number of received mails: "+strconv.Itoa(len(data)))
	}
	w.WriteHeader(200)
}

/*
	Handle the mail for a existing ticket.
 */
func (h *IncomingMailHandler) handleExistingTicketMail(ticketId int, incomingMail mail.Mail) error {
	messageEntry := ticket.MessageEntry{}
	messageEntry.Content = html.EscapeString(incomingMail.Content)
	messageEntry.CreationTime = time.Now()
	messageEntry.CreatorMail = incomingMail.Sender
	messageEntry.OnlyInternal = false
	messageEntry.CreatorMail = html.EscapeString(incomingMail.Sender)
	_, err := h.TicketContext.AppendMessageToTicket(ticketId, messageEntry)

	if err != nil {
		return err
	}
	return nil
}

/*
	Handle the mail for a new ticket.
 */
func (h *IncomingMailHandler) handleNewTicketMail( incomingMail mail.Mail) error {

	isRegistered, userId := h.UserContext.GetUserForEmail(incomingMail.Sender)
	if isRegistered {
		return h.handleNewTicketForInternalUser(userId, incomingMail)
	}
	return h.handleNewTicketForUnknownSender(incomingMail)
}

/*
	Handle the incoming mails.
 */
func (h *IncomingMailHandler) handleIncomingMails(data []mail.Mail) error {
	mailIdExtractor := newMailIdExtractor()
	for _, incomingMail := range data{
		isExistingTicket, ticketId := mailIdExtractor.getTicketId(incomingMail.Subject)
		if isExistingTicket {
			ticketReallyExists, existingTicket := h.TicketContext.GetTicketById(ticketId)
			if ticketReallyExists {
				err := h.handleExistingTicketMail(ticketId, incomingMail)
				if err != nil {
					return err
				}
				ticketCreatorMail := existingTicket.Info().Creator.Mail
				if strings.ToLower(incomingMail.Sender) != strings.ToLower(ticketCreatorMail){
					subject := "New Entry for your ticket: " + html.EscapeString(incomingMail.Subject)
					content :=  html.EscapeString(mail.BuildNotificationMailContent(existingTicket.Info().Creator.Mail,
						ticketCreatorMail, incomingMail.Content))
					err = h.MailContext.CreateNewOutgoingMail(existingTicket.Info().Creator.Mail, subject,content )
					if err != nil {
						return err
					}
				}

				continue
			} else {
				err := h.handleNewTicketMail(incomingMail)
				if err != nil {
					return err
				}
			}
		} else { // Non Existing ticket
			err := h.handleNewTicketMail(incomingMail)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
	Handle a new ticket for a internal user.
 */
func (h *IncomingMailHandler) handleNewTicketForInternalUser(userId int, incomingMail mail.Mail) error {
	exists, internalUser := h.UserContext.GetUserById(userId)
	if !exists {
		return errors.New("user should exist but does not")
	}

	title := html.EscapeString(incomingMail.Subject)
	message := ticket.MessageEntry{CreatorMail: html.EscapeString(incomingMail.Sender),
		CreationTime: time.Now(),
		Content: html.EscapeString(incomingMail.Content),
		OnlyInternal: false}
	_, err  := h.TicketContext.CreateNewTicketForInternalUser(title, internalUser, message)
	if err != nil {
		return err
	}
	return nil
}

/*
	Handle a new ticket for an unknown sender.
 */
func (h *IncomingMailHandler) handleNewTicketForUnknownSender(incomingMail mail.Mail) error {

	title := html.EscapeString(incomingMail.Subject)
	message := ticket.MessageEntry{CreatorMail: html.EscapeString(incomingMail.Sender),
		CreationTime: time.Now(),
		Content: html.EscapeString(incomingMail.Content),
		OnlyInternal: false}
	creator := ticket.Creator{Mail: html.EscapeString(incomingMail.Sender), FirstName: "SentPerMail", LastName: "SentPerMail"}
	_, err  := h.TicketContext.CreateNewTicket(title, creator, message)
	if err != nil {
		return err
	}
	return nil
}
