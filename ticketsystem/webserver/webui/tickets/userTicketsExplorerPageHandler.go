package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
	"strings"
)

/*
	Structure for the User Tickets Explorer Page handler.
*/
type UserTicketsExplorerPageHandler struct {
	TicketContext   ticket.TicketContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the User Tickets Explorer Page Data.
*/
type userTicketsExplorerPageData struct {
	Tickets []ticket.TicketInfo
	pages.BasePageData
}

/*
	The User Tickets Explorer Page handler.
*/
func (t UserTicketsExplorerPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		// Checks if the User is already logged in and if so redirects him to the start page
		isUserLoggedIn := wrappers.IsAuthenticated(r.Context())

		if !isUserLoggedIn {
			http.Redirect(w, r, "/login", http.StatusForbidden)
			return
		}

		userId := wrappers.GetUserId(r.Context())
		tickets := t.TicketContext.GetTicketsForEditorId(userId)

		data := userTicketsExplorerPageData{
			Tickets: tickets,
		}

		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = "user_tickets"

		templateRenderError := t.TemplateManager.RenderTemplate(w, "TicketExplorerPage", data)

		if templateRenderError != nil {
			t.Logger.LogError("UserTicketExplorerPage", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
