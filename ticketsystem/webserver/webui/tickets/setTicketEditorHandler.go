package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TicketSetEditorHandler struct {
	Logger        logging.Logger
	TicketContext ticketData.TicketContext
	MailContext   mailData.MailContext
	UserContext   userData.UserContext
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

		// Check ticketData:
		ticketId, err := strconv.Atoi(rawTicketId)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		ticketExists, existingTicket := t.TicketContext.GetTicketById(ticketId)
		if !ticketExists {
			t.Logger.LogError("TicketSetEditorHandler", errors.New("ticketData does not exist. id: "+rawTicketId))
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		// Check userData:
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
	existingTicket *ticketData.Ticket, rawTicketId string, ticketId int) {
	// Remove the editor:
	err := t.TicketContext.RemoveEditor(ticketId)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	// Append message for history:
	// Extract the userData who makes the change:
	loggedInUserId := wrappers.GetUserId(r.Context())
	userExists, authenticatedUser := t.UserContext.GetUserById(loggedInUserId)
	if !userExists {
		t.Logger.LogError("TicketSetEditorHandler", errors.New("userData should exist"))
		http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusBadRequest)
		return
	}
	// Build message for history:
	messageEntry := ticketData.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
		Content: "Editor removed", CreationTime: time.Now()}
	_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusInternalServerError)
		return
	}
	// Change the state if if it was not set to processing and save the change in the history:
	if existingTicket.Info().State != ticketData.Open {
		_, err = t.TicketContext.SetTicketState(ticketId, ticketData.Open)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusInternalServerError)
			return
		}
		// Build message for history:
		messageEntry := ticketData.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
			Content: "Set new state: " + ticketData.Open.String(), CreationTime: time.Now()}
		_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusInternalServerError)
			return
		}
	}

	// Notify:
	receiver := existingTicket.Info().Creator.Mail
	subject := mailData.BuildTicketEditorChangedNotificationMailSubject(ticketId)
	mailContent := mailData.BuildTicketEditorRemovedNotificationMailContent(receiver, ticketId)
	err = t.MailContext.CreateNewOutgoingMail(existingTicket.Info().Creator.Mail, subject, mailContent)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusFound)
	t.Logger.LogInfo("TicketSetEditorHandler","Editor removed from ticket: " + rawTicketId)
}

/*
	Handle a change editor request.
*/
func (t TicketSetEditorHandler) handleChangeEditorRequest(w http.ResponseWriter, r *http.Request, existingTicket *ticketData.Ticket, rawTicketId string, ticketId int, userId int) {
	userExists, existingUser := t.UserContext.GetUserById(userId)
	if !userExists {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("userData does not exist"))
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Get the userData who makes the change:
	loggedInUserId := wrappers.GetUserId(r.Context())
	userExists, authenticatedUser := t.UserContext.GetUserById(loggedInUserId)
	if !userExists {
		t.Logger.LogError("TicketSetEditorHandler", errors.New("userData should exist"))
		http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusBadRequest)
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
	messageEntry := ticketData.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
		Content: "Set new editor: " + existingUser.GetUserNameString(), CreationTime: time.Now()}
	_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusInternalServerError)
		return
	}
	// Change the state if if it was not set to processing and save the change in the history:
	if existingTicket.Info().State != ticketData.Processing {
		_, err = t.TicketContext.SetTicketState(ticketId, ticketData.Processing)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusInternalServerError)
			return
		}
		// Write the state change to the history:
		messageEntry := ticketData.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
			Content: "Set new state: " + ticketData.Processing.String(), CreationTime: time.Now()}
		_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
		if err != nil {
			t.Logger.LogError("TicketSetEditorHandler", err)
			http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusInternalServerError)
			return
		}
	}

	// Notify:
	receiver := existingTicket.Info().Creator.Mail
	subject := mailData.BuildTicketEditorChangedNotificationMailSubject(ticketId)
	mailContent := mailData.BuildTicketEditorChangedNotificationMailContent(receiver, ticketId, existingUser.GetUserNameString())
	err = t.MailContext.CreateNewOutgoingMail(existingTicket.Info().Creator.Mail, subject, mailContent)
	if err != nil {
		t.Logger.LogError("TicketSetEditorHandler", err)
		http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/ticketData/"+rawTicketId, http.StatusFound)
	t.Logger.LogInfo("TicketSetEditorHandler","Editor set for ticket: " + rawTicketId + ". New editorId: " + strconv.Itoa(existingUser.UserId))
}
