// 5894619, 6720876, 9793350
package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
	"html"
	"net/http"
	"strconv"
	"strings"
)

/*
	Structure for the Tickets View Page handler.
*/
type TicketViewPageHandler struct {
	UserContext     userData.UserContext
	TicketContext   ticketData.TicketContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the Ticket View Page Data.
*/
type ticketViewPageData struct {
	TicketInfo ticketData.TicketInfo
	Messages   []ticketData.MessageEntry
	UserName   string
	pages.BasePageData
}

/*
	The Ticket View Page handler.
*/
func (t TicketViewPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		isUserLoggedIn := wrappers.IsAuthenticated(r.Context())

		URLPath := strings.Split(r.URL.Path, "/")

		urlValue := html.EscapeString(URLPath[2])

		ticketId, idConversionError := strconv.Atoi(urlValue)

		if idConversionError != nil {
			t.Logger.LogError("TicketViewPageHandler", idConversionError)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		ticketExist, existingTicket := t.TicketContext.GetTicketById(ticketId)

		if !ticketExist {
			t.Logger.LogError("TicketViewPageHandler", errors.New("Ticket doesn´t exist!"))
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}

		ticketInfo := existingTicket.Info()
		messages := existingTicket.Messages()
		mail := ""

		if isUserLoggedIn {
			userId := wrappers.GetUserId(r.Context())
			exists, existingUser := t.UserContext.GetUserById(userId)

			if exists {
				mail = existingUser.Mail
			}
		} else {
			messages = filterOutInternalOnlyMessages(messages)
		}

		data := ticketViewPageData{
			TicketInfo: ticketInfo,
			Messages:   messages,
			UserName:   mail,
		}

		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = "all_tickets"

		templateRenderError := t.TemplateManager.RenderTemplate(w, "TicketViewPage", data)

		if templateRenderError != nil {
			t.Logger.LogError("TicketViewPageHandler", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}

func filterOutInternalOnlyMessages(messages []ticketData.MessageEntry) (externalMessages []ticketData.MessageEntry) {
	for _, message := range messages {
		if !message.OnlyInternal {
			externalMessages = append(externalMessages, message)
		}
	}
	return
}
