package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

/*
	Structure for the Ticket Create Page handler.
*/
type TicketCreatePageHandler struct {
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
	UserContext     userData.UserContext
}

/*
	Structure for the Ticket Create Page Data.
*/
type ticketCreatePageData struct {
	pages.BasePageData
	IsUserLoggedIn bool
	UserName       string
	FirstName      string
	LastName       string
}

/*
	The Ticket Create Page handler.
*/
func (t TicketCreatePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		isUserLoggedIn := wrappers.IsAuthenticated(r.Context())

		data := ticketCreatePageData{
			UserName:       "",
			IsUserLoggedIn: false,
			FirstName:      "",
			LastName:       "",
		}

		if isUserLoggedIn {
			userId := wrappers.GetUserId(r.Context())

			userExist, existingUser := t.UserContext.GetUserById(userId)

			if userExist {
				data = ticketCreatePageData{
					UserName:       existingUser.Mail,
					IsUserLoggedIn: isUserLoggedIn,
					FirstName:      existingUser.FirstName,
					LastName:       existingUser.LastName,
				}
			} else {
				t.Logger.LogError("TicketCreatePageHandler", errors.New("User ID couldn´t be referenced back to a user!"))
				http.Redirect(w, r, "/", http.StatusInternalServerError)
				return
			}
		}

		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = "ticket_create"

		templateRenderError := t.TemplateManager.RenderTemplate(w, "TicketCreatePage", data)

		if templateRenderError != nil {
			t.Logger.LogError("TicketCreatePageHandler", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}

}
