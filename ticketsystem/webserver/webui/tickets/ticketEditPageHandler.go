package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"errors"
	"html"
	"net/http"
	"strconv"
	"strings"
)

/*
	Structure for the Tickets Edit Page handler.
*/
type TicketEditPageHandler struct {
	TicketContext   ticket.TicketContext
	UserContext		user.UserContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the Ticket Edit Page Data.
*/
type ticketEditPageData struct {
	TicketInfo 		ticket.TicketInfo
	OtherTickets	[]ticket.TicketInfo
	Users			[]user.User
	pages.BasePageData
}

/*
	The Ticket Edit Page handler.
*/
func (t TicketEditPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		URLPath := strings.Split(r.URL.Path, "/")

		urlValue := html.EscapeString(URLPath[2])

		ticketId, idConversionError := strconv.Atoi(urlValue)

		if idConversionError != nil {
			t.Logger.LogError("TicketEditPageHandler", idConversionError)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		ticketExist, editableTicket := t.TicketContext.GetTicketById(ticketId)

		if !ticketExist {
			t.Logger.LogError("TicketEditPageHandler", errors.New("Ticket doesn´t exist."))
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}

		ticketInfo := editableTicket.Info()

		users := t.UserContext.GetAllActiveUsers()

		otherTickets := t.TicketContext.GetAllTicketInfo()

		otherTickets = filterOutTicket(ticketInfo.Id, otherTickets)

		data := ticketEditPageData{
			TicketInfo: 	ticketInfo,
			OtherTickets:	otherTickets,
			Users:			users,
		}

		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = "ticket_edit"

		templateRenderError := t.TemplateManager.RenderTemplate(w, "TicketEditPage", data)

		if templateRenderError != nil {
			t.Logger.LogError("TicketEditPageHandler", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}

func filterOutTicket(ticketIdToFilter int, tickets []ticket.TicketInfo) (result []ticket.TicketInfo) {
	for _, t := range tickets {
		if t.Id != ticketIdToFilter {
			result = append(result, t)
		}
	}
	return
}