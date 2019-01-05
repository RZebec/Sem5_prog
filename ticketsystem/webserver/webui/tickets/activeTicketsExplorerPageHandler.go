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
	Structure for the active Ticket Explorer Page handler.
*/
type ActiveTicketsExplorerPageHandler struct {
	TicketContext   ticketData.TicketContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the active Ticket Explorer Page Data.
*/
type activeTicketsExplorerPageData struct {
	Tickets []ticketData.TicketInfo
	pages.BasePageData
}

/*
	The active Ticket Explorer Page handler.
*/
func (t ActiveTicketsExplorerPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		existingTickets := t.TicketContext.GetAllTicketInfo()
		var filteredTickets []ticketData.TicketInfo
		for _, curTicket := range existingTickets {
			if curTicket.State == ticketData.Processing {
				filteredTickets = append(filteredTickets, curTicket)
			}
		}

		sort.Slice(filteredTickets, func(i, j int) bool {
			return filteredTickets[i].CreationTime.Before(filteredTickets[j].CreationTime)
		})

		data := activeTicketsExplorerPageData{
			Tickets: filteredTickets,
		}

		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = "active_tickets"

		templateRenderError := t.TemplateManager.RenderTemplate(w, "TicketExplorerPage", data)

		if templateRenderError != nil {
			t.Logger.LogError("ActiveTicketsExplorerPageHandler", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
