package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"net/http"
)

/*
	Structure for the Login handler.
*/
type AdminPageHandler struct {
	UserContext user.UserContext
	Logger      logging.Logger
	TemplateManager	templateManager.TemplateContext
	ApiContext	config.ApiContext
}

/*
	Structure for the Login Page Data.
*/
type adminPageData struct {
	Users 		[]user.User
	IncomingMailApiKey	string
	OutgoingMailApiKey	string
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

		templateRenderError := a.TemplateManager.RenderTemplate(w, "AdminPage", data)

		if templateRenderError != nil {
			a.Logger.LogError("Admin", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
