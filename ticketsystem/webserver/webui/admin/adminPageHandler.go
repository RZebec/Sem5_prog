package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
)

/*
	Structure for the Admin handler.
*/
type AdminPageHandler struct {
	UserContext     user.UserContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
	ApiContext      config.ApiContext
}

/*
	Structure for the Admin Page Data.
*/
type adminPageData struct {
	Users              []user.User
	IncomingMailApiKey string
	OutgoingMailApiKey string
	pages.BasePageData
}

/*
	The Admin Page handler.
*/
func (a AdminPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		users := a.UserContext.GetAllLockedUsers()

		incomingMailApiKey := a.ApiContext.GetIncomingMailApiKey()
		outgoingMailApiKey := a.ApiContext.GetOutgoingMailApiKey()

		data := adminPageData{
			Users:              users,
			IncomingMailApiKey: incomingMailApiKey,
			OutgoingMailApiKey: outgoingMailApiKey,
		}
		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = "admin"

		templateRenderError := a.TemplateManager.RenderTemplate(w, "AdminPage", data)

		if templateRenderError != nil {
			a.Logger.LogError("Admin", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
