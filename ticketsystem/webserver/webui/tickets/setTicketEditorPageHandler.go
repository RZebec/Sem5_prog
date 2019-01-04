package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TicketSetEditorHandler struct {
	Logger        logging.Logger
	TicketContext ticket.TicketContext
	MailContext   mail.MailContext
	UserContext   user.UserContext
}

/*
	Handling a change editor request.
*/
func (t TicketSetEditorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		rawUserId := r.FormValue("editorUserId")
		rawTicketId := r.FormValue("ticketId")

		// Check ticket:
		ticketId, err := strconv.Atoi(rawTicketId)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		ticketExists, existingTicket := t.TicketContext.GetTicketById(ticketId)
		if !ticketExists {
			t.Logger.LogError("TicketSetEditorHandler", errors.New("ticket does not exist. id: "+rawTicketId))
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		// Check user:
		userId, err := strconv.Atoi(rawUserId)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		if userId == -1 {
			t.handleRemoveEditorRequest(w, r, existingTicket, rawTicketId, ticketId)
			return
		}
		t.handleChangeEditorRequest(w, r, existingTicket, rawTicketId, ticketId, userId)
	}
}

/*
	Handle remove editor requests.
*/
func (t TicketSetEditorHandler) handleRemoveEditorRequest(w http.ResponseWriter, r *http.Request,
	existingTicket *ticket.Ticket, rawTicketId string, ticketId int) {
	// Remove the editor:
	err := t.TicketContext.RemoveEditor(ticketId)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	// Append message for history:
	// Extract the user who makes the change:
	loggedInUserId := wrappers.GetUserId(r.Context())
	userExists, authenticatedUser := t.UserContext.GetUserById(loggedInUserId)
	if !userExists {
		t.Logger.LogError("TicketSetEditorHandler", errors.New("user should exist"))
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusBadRequest)
		return
	}
	// Build message for history:
	messageEntry := ticket.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
		Content: "Editor removed", CreationTime: time.Now()}
	_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
		return
	}
	// Change the state if if it was not set to processing and save the change in the history:
	if existingTicket.Info().State != ticket.Processing {
		_, err = t.TicketContext.SetTicketState(ticketId, ticket.Processing)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
			return
		}
		// Build message for history:
		messageEntry := ticket.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
			Content: "Set new state: " + ticket.Processing.String(), CreationTime: time.Now()}
		_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
			return
		}
	}

	// Notify:
	receiver := existingTicket.Info().Creator.Mail
	subject := mail.BuildTicketEditorChangedNotificationMailSubject(ticketId)
	mailContent := mail.BuildTicketEditorRemovedNotificationMailContent(receiver, ticketId)
	err = t.MailContext.CreateNewOutgoingMail(existingTicket.Info().Creator.Mail, subject, mailContent)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusOK)
}

/*
	Handle a change editor request.
*/
func (t TicketSetEditorHandler) handleChangeEditorRequest(w http.ResponseWriter, r *http.Request, existingTicket *ticket.Ticket, rawTicketId string, ticketId int, userId int) {
	userExists, existingUser := t.UserContext.GetUserById(userId)
	if !userExists {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("user does not exist"))
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Get the user who makes the change:
	loggedInUserId := wrappers.GetUserId(r.Context())
	userExists, authenticatedUser := t.UserContext.GetUserById(loggedInUserId)
	if !userExists {
		t.Logger.LogError("TicketSetEditorHandler", errors.New("user should exist"))
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusBadRequest)
		return
	}

	// Set editor:
	_, err := t.TicketContext.SetEditor(existingUser, ticketId)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	// Append message for history:
	// Build message for history:
	messageEntry := ticket.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
		Content: "Set new editor: " + existingUser.GetUserNameString(), CreationTime: time.Now()}
	_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
		return
	}
	// Change the state if if it was not set to processing and save the change in the history:
	if existingTicket.Info().State != ticket.Processing {
		_, err = t.TicketContext.SetTicketState(ticketId, ticket.Processing)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
			return
		}
		// Write the state change to the history:
		messageEntry := ticket.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
			Content: "Set new state: " + ticket.Processing.String(), CreationTime: time.Now()}
		_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
			return
		}
	}

	// Notify:
	receiver := existingTicket.Info().Creator.Mail
	subject := mail.BuildTicketEditorChangedNotificationMailSubject(ticketId)
	mailContent := mail.BuildTicketEditorChangedNotificationMailContent(receiver, ticketId, existingUser.GetUserNameString())
	err = t.MailContext.CreateNewOutgoingMail(existingTicket.Info().Creator.Mail, subject, mailContent)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusOK)
}
