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
	Structure for the Closed Ticket Explorer Page handler.
*/
type ClosedTicketsExplorerPageHandler struct {
	TicketContext   ticketData.TicketContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the Closed Ticket Explorer Page Data.
*/
type closedTicketsExplorerPageData struct {
	Tickets []ticketData.TicketInfo
	pages.BasePageData
}

/*
	The Closed Ticket Explorer Page handler.
*/
func (t ClosedTicketsExplorerPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		existingTickets := t.TicketContext.GetAllTicketInfo()

		var tickets []ticketData.TicketInfo

		for _, curTicket := range existingTickets {
			if curTicket.State == ticketData.Closed {
				tickets = append(tickets, curTicket)
			}
		}

		sort.Slice(tickets, func(i, j int) bool {
			return tickets[i].CreationTime.Before(tickets[j].CreationTime)
		})

		data := closedTicketsExplorerPageData{
			Tickets: tickets,
		}

		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = "closed_tickets"

		templateRenderError := t.TemplateManager.RenderTemplate(w, "TicketExplorerPage", data)

		if templateRenderError != nil {
			t.Logger.LogError("AllTicketExplorerPage", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
