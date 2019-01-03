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
	Structure for the Open Tickets Explorer Page handler.
*/
type OpenTicketsExplorerPageHandler struct {
	TicketContext   ticket.TicketContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the Open Tickets Explorer Page Data.
*/
type openTicketsExplorerPageData struct {
	Tickets []ticket.TicketInfo
	pages.BasePageData
}

/*
	The Open Tickets Explorer Page handler.
*/
func (t OpenTicketsExplorerPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		tickets := t.TicketContext.GetAllOpenTickets()

		data := openTicketsExplorerPageData{
			Tickets: tickets,
		}

		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = "open_tickets"

		templateRenderError := t.TemplateManager.RenderTemplate(w, "TicketExplorerPage", data)

		if templateRenderError != nil {
			t.Logger.LogError("OpenTicketExplorerPage", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
