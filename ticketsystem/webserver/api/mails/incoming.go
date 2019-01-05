package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
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
	Logger            logging.Logger
	MailContext       mailData.MailContext
	TicketContext     ticketData.TicketContext
	UserContext       userData.UserContext
	MailRepliesFilter AutomaticRepliesFilter
}

/*
	Handling the incoming mails.
*/
func (h *IncomingMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data []mailData.Mail
	err := decoder.Decode(&data)
	if err != nil {
		h.Logger.LogError("IncomingMailHandler", err)
		w.WriteHeader(400)
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
	Handle the mail for a existing ticketData.
*/
func (h *IncomingMailHandler) handleExistingTicketMail(ticketId int, incomingMail mailData.Mail) error {
	_, err := h.TicketContext.AppendMessageToTicket(ticketId, h.buildMessageEntry(incomingMail))
	if err != nil {
		return err
	}
	return nil
}

/*
	Build a message entry.
*/
func (h *IncomingMailHandler) buildMessageEntry(incomingMail mailData.Mail) ticketData.MessageEntry {
	messageEntry := ticketData.MessageEntry{}
	messageEntry.Content = html.EscapeString(incomingMail.Content)
	messageEntry.CreationTime = time.Now()
	messageEntry.CreatorMail = incomingMail.Sender
	messageEntry.OnlyInternal = false
	messageEntry.CreatorMail = html.EscapeString(incomingMail.Sender)
	return messageEntry
}

/*
	Handle the mail for a new ticketData.
*/
func (h *IncomingMailHandler) handleNewTicketMail(incomingMail mailData.Mail) error {
	isRegistered, userId := h.UserContext.GetUserForEmail(incomingMail.Sender)
	if isRegistered {
		return h.handleNewTicketForInternalUser(userId, incomingMail)
	}
	return h.handleNewTicketForUnknownSender(incomingMail)
}

/*
	Handle the incoming mails.
*/
func (h *IncomingMailHandler) handleIncomingMails(data []mailData.Mail) error {
	mailIdExtractor := newMailIdExtractor()
	for _, incomingMail := range data {
		if h.MailRepliesFilter.IsAutomaticResponse(incomingMail) {
			h.Logger.LogWarning("IncomingMailHandler", "Ignoring mail with id "+incomingMail.Id+" because"+
				"it is a automatic reply")
			continue
		}
		isExistingTicket, ticketId := mailIdExtractor.getTicketId(incomingMail.Subject)
		if isExistingTicket {
			ticketReallyExists, existingTicket := h.TicketContext.GetTicketById(ticketId)
			if ticketReallyExists {
				err := h.handleExistingTicketMail(ticketId, incomingMail)
				if err != nil {
					return err
				}
				ticketCreatorMail := existingTicket.Info().Creator.Mail
				if strings.ToLower(incomingMail.Sender) != strings.ToLower(ticketCreatorMail) {
					subject := "New Entry for your ticketData: " + html.EscapeString(incomingMail.Subject)
					senderOfMail := html.EscapeString(incomingMail.Sender)
					content := html.EscapeString(mailData.BuildAppendMessageNotificationMailContent(ticketCreatorMail,
						senderOfMail, incomingMail.Content))
					err = h.MailContext.CreateNewOutgoingMail(existingTicket.Info().Creator.Mail, subject, content)
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
		} else { // Non Existing ticketData
			err := h.handleNewTicketMail(incomingMail)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
	Handle a new ticketData for a internal userData.
*/
func (h *IncomingMailHandler) handleNewTicketForInternalUser(userId int, incomingMail mailData.Mail) error {
	exists, internalUser := h.UserContext.GetUserById(userId)
	if !exists {
		return errors.New("userData should exist but does not")
	}

	title := html.EscapeString(incomingMail.Subject)
	message := h.buildMessageEntry(incomingMail)
	_, err := h.TicketContext.CreateNewTicketForInternalUser(title, internalUser, message)
	if err != nil {
		return err
	}
	return nil
}

/*
	Handle a new ticketData for an unknown sender.
*/
func (h *IncomingMailHandler) handleNewTicketForUnknownSender(incomingMail mailData.Mail) error {
	title := html.EscapeString(incomingMail.Subject)
	message := h.buildMessageEntry(incomingMail)
	creator := h.buildCreator(incomingMail)
	_, err := h.TicketContext.CreateNewTicket(title, creator, message)
	if err != nil {
		return err
	}
	return nil
}

/*
	Build the creator.
*/
func (h *IncomingMailHandler) buildCreator(incomingMail mailData.Mail) ticketData.Creator {
	return ticketData.Creator{Mail: html.EscapeString(incomingMail.Sender), FirstName: "SentPerMail", LastName: "SentPerMail"}
}
