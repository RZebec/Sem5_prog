package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"github.com/pkg/errors"
	"html"
	"net/http"
	"strconv"
	"strings"
)

/*
	Structure for the Ticket Create handler.
*/
type TicketCreateHandler struct {
	UserContext     user.UserContext
	Logger          logging.Logger
	TicketContext	ticket.TicketContext
}

/*
	The Ticket Create handler.
*/
func (t TicketCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		mail := r.FormValue("mail")
		title := r.FormValue("title")
		message := r.FormValue("message")
		internal := r.FormValue("internal")

		mail = html.EscapeString(mail)
		title = html.EscapeString(title)
		message = html.EscapeString(message)
		internal = html.EscapeString(internal)

		internalOnly, parseError := strconv.ParseBool(internal)

		if parseError != nil {
			t.Logger.LogError("TicketCreateHandler", parseError)
			http.Redirect(w, r, "/ticket_create", http.StatusInternalServerError)
			return
		}

		exist, userId := t.UserContext.GetUserForEmail(mail)

		initialMessage := ticket.MessageEntry{Id: 0, CreatorMail: mail, Content: message, OnlyInternal: internalOnly}

		if !exist {
			ticket, err := t.TicketContext.CreateNewTicket(title, ticket.Creator{Mail: mail}, initialMessage)

			if err != nil {
				t.Logger.LogError("TicketCreateHandler", err)
				http.Redirect(w, r, "/ticket_create", http.StatusInternalServerError)
				return
			}

			ticketId := strconv.Itoa(ticket.Info().Id)

			http.Redirect(w, r, "/ticket/" + ticketId, 302)
			return
		}

		exist, user := t.UserContext.GetUserById(userId)

		if (!exist) {
			t.Logger.LogError("TicketCreateHandler", errors.New("User doesn´t exist."))
			http.Redirect(w, r, "/ticket_create", http.StatusInternalServerError)
			return
		}

		ticket, err := t.TicketContext.CreateNewTicketForInternalUser(title, user, initialMessage)

		if err != nil {
			t.Logger.LogError("TicketCreateHandler", err)
			http.Redirect(w, r, "/ticket_create", http.StatusInternalServerError)
			return
		}

		ticketId := strconv.Itoa(ticket.Info().Id)

		http.Redirect(w, r, "/ticket/" + ticketId, 302)
	}
}
