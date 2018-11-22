package ticketView

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"net/http"
	"strconv"
	"strings"
)

/*
	Structure for the Ticket Explorer Page Handler.
*/
type TicketViewPageHandler struct {
	Config      config.Configuration
	UserContext user.UserContext
	TicketContext ticket.TicketContext
}

/*
	Structure for the Ticket Explorer Page Data.
*/
type TicketViewPageData struct {
	TicketInfo ticket.TicketInfo
	Messages []ticket.MessageEntry
	IsUserLoggedIn bool
}

/*
	The Ticket Explorer Page handler.
*/
func (t TicketViewPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Checks if the User is logged in
	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, t.UserContext, t.Config.AccessTokenCookieName)

	if isUserLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	URLPath := strings.Split(r.URL.Path, "/")

	id, idConversionError := strconv.Atoi(URLPath[2])

	if idConversionError != nil {
		//TODO: Handle Error
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	ticket, err := t.TicketContext.GetTicketById(id)

	if err != nil {
		//TODO: Handle error
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	data := TicketViewPageData {
		IsUserLoggedIn: isUserLoggedIn,
		TicketInfo: ticket.Info(),
		Messages: ticket.Messages(),
	}

	templateRenderError := templateManager.RenderTemplate(w, "TicketViewPage", data)

	if templateRenderError != nil {
		// TODO: Handle error
		return
	}
}
