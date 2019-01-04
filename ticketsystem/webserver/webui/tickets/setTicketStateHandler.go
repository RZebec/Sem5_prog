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

/*
	A handler to change the state of a ticket.
*/
type SetTicketStateHandler struct {
	Logger        logging.Logger
	TicketContext ticket.TicketContext
	MailContext   mail.MailContext
	UserContext   user.UserContext
}

/*
	Resolve the state from a string.
*/
func (t *SetTicketStateHandler) ResolveState(state string) (valid bool, ticketState ticket.TicketState) {
	if strings.ToLower(state) == strings.ToLower(ticket.Open.String()) {
		return true, ticket.Open
	}
	if strings.ToLower(state) == strings.ToLower(ticket.Processing.String()) {
		return true, ticket.Processing
	}
	if strings.ToLower(state) == strings.ToLower(ticket.Closed.String()) {
		return true, ticket.Closed
	}
	return false, ticket.Open
}

/*
	Handling a change state request.
*/
func (t SetTicketStateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		rawTicketId := r.FormValue("ticketId")

		// Check ticket:
		ticketId, err := strconv.Atoi(rawTicketId)
		if err != nil {
			t.Logger.LogError("SetTicketStateHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		ticketExists, existingTicket := t.TicketContext.GetTicketById(ticketId)
		if !ticketExists {
			t.Logger.LogError("SetTicketStateHandler", errors.New("ticket does not exist. id: "+rawTicketId))
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		// Extract the user who makes the change:
		loggedInUserId := wrappers.GetUserId(r.Context())
		userExists, authenticatedUser := t.UserContext.GetUserById(loggedInUserId)
		if !userExists {
			t.Logger.LogError("SetTicketStateHandler", errors.New("user should exist"))
			http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusBadRequest)
			return
		}

		rawState := r.FormValue("newState")
		stateIsValid, newState := t.ResolveState(rawState)
		if !stateIsValid {
			t.Logger.LogError("SetTicketStateHandler", errors.New("state is not valid"))
			http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusBadRequest)
			return
		}

		// Change the state, when it is different to the existing one.
		if existingTicket.Info().State != newState {
			_, err = t.TicketContext.SetTicketState(ticketId, newState)
			if err != nil {
				t.Logger.LogError("SetTicketStateHandler", err)
				http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
				return
			}
			// Build message for history:
			messageEntry := ticket.MessageEntry{CreatorMail: authenticatedUser.Mail, OnlyInternal: false,
				Content: "Set new state: " + ticket.Processing.String(), CreationTime: time.Now()}
			_, err = t.TicketContext.AppendMessageToTicket(ticketId, messageEntry)
			if err != nil {
				t.Logger.LogError("SetTicketStateHandler", err)
				http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/ticket/"+rawTicketId, http.StatusOK)
	}
}
