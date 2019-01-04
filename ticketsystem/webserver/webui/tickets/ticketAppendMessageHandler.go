package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/validation/mail"
	mailContext "de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"errors"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
	A handler to append messages to a ticket.
*/
type TicketAppendMessageHandler struct {
	UserContext   user.UserContext
	Logger        logging.Logger
	TicketContext ticket.TicketContext
	MailContext	  mailContext.MailContext
}

/*
	Append a message to a ticket
*/
func (t TicketAppendMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		userIsAuthenticated := wrappers.IsAuthenticated(r.Context())
		if userIsAuthenticated {
			t.handlerForAuthenticatedUser(w, r)
			return
		} else {
			t.handlerForNonAuthenticatedUser(w, r)
			return
		}
	}
}

/*
	Handle message append for authenticated user.
*/
func (t TicketAppendMessageHandler) handlerForAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	// Handle ticket id:
	rawTicketId := r.FormValue("ticketId")
	ticketId, err := strconv.Atoi(rawTicketId)
	if err != nil {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("invalid ticket id"))
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	tickedExists, existingTicket := t.TicketContext.GetTicketById(ticketId)
	if !tickedExists {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("invalid ticket id, ticket does not exist"))
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Handle user:
	userId := wrappers.GetUserId(r.Context())
	userExists, authenticatedUser := t.UserContext.GetUserById(userId)
	if !userExists {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("user should exist"))
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusBadRequest)
		return
	}

	// Handle content:
	content := r.FormValue("messageContent")
	content = html.EscapeString(content)
	if len(content) == 0 {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("message content to append is empty"))
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusBadRequest)
		return
	}

	// Handle onlyInternal flag
	rawOnlyInternal := r.FormValue("onlyInternal")
	onlyInternal, err := strconv.ParseBool(rawOnlyInternal)
	if err != nil {
		t.Logger.LogError("TicketAppendMessageHandler", err)
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusBadRequest)
		return
	}

	// Append message to ticket:
	messageEntry := ticket.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: onlyInternal,
		Content: content, CreationTime: time.Now()}

	_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
	if err != nil {
		t.Logger.LogError("TicketAppendMessageHandler", err)
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
		return
	}

	// Notify the creator:
	sender := strings.ToLower(authenticatedUser.Mail)
	receiver := strings.ToLower(existingTicket.Info().Creator.Mail)
	if sender != receiver {
		subject := mailContext.BuildAppendMessageNotificationMailSubject(ticketId)
		mailContent := mailContext.BuildAppendMessageNotificationMailContent(receiver, sender, content)
		err = t.MailContext.CreateNewOutgoingMail(existingTicket.Info().Creator.Mail, subject, mailContent)
		if err != nil {
			t.Logger.LogError("TicketAppendMessageHandler", err)
			http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusFound)
}

/*
	Handle message for non authenticated user.
*/
func (t TicketAppendMessageHandler) handlerForNonAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	// Handle ticket id:
	rawTicketId := r.FormValue("ticketId")
	ticketId, err := strconv.Atoi(rawTicketId)
	if err != nil {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("invalid ticket id"))
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	tickedExists, existingTicket := t.TicketContext.GetTicketById(ticketId)
	if !tickedExists {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("invalid ticket id. ticket does not exist."))
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Handle mail:
	rawMail := r.FormValue("mail")
	rawMail = html.EscapeString(rawMail)

	validator := mail.NewValidator()
	mailIsValid := validator.Validate(rawMail)

	if !mailIsValid {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("mail is invalid"))
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusBadRequest)
		return
	}

	// It should not be possible to append a message for a user which is registered, if the current user is not logged in.
	userExists, _ := t.UserContext.GetUserForEmail(rawMail)
	if userExists {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("user is registered but not logged in"))
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Handle content:
	content := r.FormValue("messageContent")
	content = html.EscapeString(content)
	if len(content) == 0 {
		t.Logger.LogError("TicketAppendMessageHandler", errors.New("message content to append is empty"))
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusBadRequest)
		return
	}

	// Append message:
	messageEntry := ticket.MessageEntry{CreatorMail: rawMail, OnlyInternal: false,
		Content: content, CreationTime: time.Now()}

	_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
	if err != nil {
		t.Logger.LogError("TicketAppendMessageHandler", err)
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
		return
	}

	// Notify the creator:
	sender := strings.ToLower(rawMail)
	receiver := strings.ToLower(existingTicket.Info().Creator.Mail)
	if sender != receiver {
		subject := mailContext.BuildAppendMessageNotificationMailSubject(ticketId)
		mailContent := mailContext.BuildAppendMessageNotificationMailContent(receiver, sender, content)
		err = t.MailContext.CreateNewOutgoingMail(existingTicket.Info().Creator.Mail, subject, mailContent)
		if err != nil {
			t.Logger.LogError("TicketAppendMessageHandler", err)
			http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusFound)
}
