package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
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
	TicketContext   ticketData.TicketContext
	UserContext     userData.UserContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the Ticket Edit Page Data.
*/
type ticketEditPageData struct {
	TicketInfo   ticketData.TicketInfo
	OtherTickets []ticketData.TicketInfo
	Users        []userData.User
	OtherState1  ticketData.TicketState
	OtherState2  ticketData.TicketState
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

		urlValue := html.EscapeString(URLPath[3])

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

		otherTickets = filterOutTicket(ticketInfo, otherTickets)

		otherStates := getOtherTicketStates(ticketInfo.State)

		data := ticketEditPageData{
			TicketInfo: 	ticketInfo,
			OtherTickets:	otherTickets,
			Users:			users,
			OtherState1:	otherStates[0],
			OtherState2:	otherStates[1],
		}
		
		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = ""

		templateRenderError := t.TemplateManager.RenderTemplate(w, "TicketEditPage", data)

		if templateRenderError != nil {
			t.Logger.LogError("TicketEditPageHandler", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}

func filterOutTicket(ticketInfo ticketData.TicketInfo, tickets []ticketData.TicketInfo) (result []ticketData.TicketInfo) {
	for _, t := range tickets {
		if t.Id != ticketInfo.Id && t.Editor == ticketInfo.Editor {
			result = append(result, t)
		}
	}
	return
}

func getOtherTicketStates(ticketState ticketData.TicketState) (result []ticketData.TicketState) {
	states := []ticketData.TicketState{ticketData.Open, ticketData.Closed, ticketData.Processing}
	for _, state := range states {
		if state != ticketState {
			result = append(result, state)
		}
	}
	return
}
