package ticketExplorer

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"net/http"
)

/*
	Structure for the Ticket Explorer Page Handler.
*/
type TicketExplorerPageHandler struct {
	Config      config.Configuration
	UserContext user.UserContext
	TicketContext ticket.TicketContext
}

/*
	Structure for the Ticket Explorer Page Data.
*/
type TicketExplorerPageData struct {
	TicketInfo []ticket.TicketInfo
	IsUserLoggedIn bool
}

/*
	The Ticket Explorer Page handler.
*/
func (t TicketExplorerPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Checks if the User is logged in
	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, t.UserContext, t.Config.AccessTokenCookieName)

	if isUserLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	ticketInfo := t.TicketContext.GetAllTicketInfo()

	data := TicketExplorerPageData {
		IsUserLoggedIn: isUserLoggedIn,
		TicketInfo: ticketInfo,
	}

	templateRenderError := templateManager.RenderTemplate(w, "TicketExplorerPage", data)

	if templateRenderError != nil {
		// TODO: Handle error
	}
}