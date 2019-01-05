package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
	"sort"
	"strings"
)

/*
	Structure for the User Tickets Explorer Page handler.
*/
type UserTicketsExplorerPageHandler struct {
	TicketContext   ticketData.TicketContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the User Tickets Explorer Page Data.
*/
type userTicketsExplorerPageData struct {
	Tickets []ticketData.TicketInfo
	pages.BasePageData
}

/*
	The User Tickets Explorer Page handler.
*/
func (t UserTicketsExplorerPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		userId := wrappers.GetUserId(r.Context())
		tickets := t.TicketContext.GetTicketsForEditorId(userId)

		sort.Slice(tickets, func(i, j int) bool {
			return tickets[i].CreationTime.Before(tickets[j].CreationTime)
		})

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
