package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
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
	UserContext   userData.UserContext
	Logger        logging.Logger
	TicketContext ticketData.TicketContext
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
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		message := r.FormValue("message")
		internal := r.FormValue("internal")

		mail = html.EscapeString(mail)
		title = html.EscapeString(title)
		firstName = html.EscapeString(firstName)
		lastName = html.EscapeString(lastName)
		message = html.EscapeString(message)
		internal = html.EscapeString(internal)

		internalOnly, parseError := strconv.ParseBool(internal)

		if internal == ""  {
			internalOnly = false
		} else if parseError != nil{
			t.Logger.LogError("TicketCreateHandler", parseError)
			http.Redirect(w, r, "/ticket_create", http.StatusInternalServerError)
			return
		}

		isUserLoggedIn := wrappers.IsAuthenticated(r.Context())

		loggedInUserId := wrappers.GetUserId(r.Context())

		exist, userId := t.UserContext.GetUserForEmail(mail)

		if loggedInUserId != userId {
			t.Logger.LogError("TicketCreateHandler", errors.New("User with the corresponding mail address is not logged in!"))
			http.Redirect(w, r, "/ticket_create", http.StatusBadRequest)
			return
		}

		if !isUserLoggedIn && !exist {
			initialMessage := ticketData.MessageEntry{Id: 0, CreatorMail: mail, Content: message, OnlyInternal: false}

			createdTicket, err := t.TicketContext.CreateNewTicket(title, ticketData.Creator{Mail: mail, FirstName: firstName, LastName: lastName}, initialMessage)

			if err != nil {
				t.Logger.LogError("TicketCreateHandler", err)
				http.Redirect(w, r, "/ticket_create", http.StatusInternalServerError)
				return
			}

			ticketId := strconv.Itoa(createdTicket.Info().Id)

			http.Redirect(w, r, "/ticket/" + ticketId, 302)
			return
		}

		exist, existingUser := t.UserContext.GetUserById(userId)

		if !exist {
			t.Logger.LogError("TicketCreateHandler", errors.New("User doesn´t exist."))
			http.Redirect(w, r, "/ticket_create", http.StatusInternalServerError)
			return
		}

		initialMessage := ticketData.MessageEntry{Id: 0, CreatorMail: mail, Content: message, OnlyInternal: internalOnly}

		createdTicket, err := t.TicketContext.CreateNewTicketForInternalUser(title, existingUser, initialMessage)

		if err != nil {
			t.Logger.LogError("TicketCreateHandler", err)
			http.Redirect(w, r, "/ticket_create", http.StatusInternalServerError)
			return
		}

		ticketId := strconv.Itoa(createdTicket.Info().Id)

		http.Redirect(w, r, "/ticket/" + ticketId, 302)
		t.Logger.LogInfo("TicketCreateHandler", "New ticket with id " + ticketId + "created")
	}
}
