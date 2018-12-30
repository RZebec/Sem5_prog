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
	Config      config.Configuration
	Logger      logging.Logger
}

/*
	Structure for the Login Page Data.
*/
type adminPageData struct {
	Users 		[]user.User
}

/*
	The Admin Page handler.
*/
func (a AdminPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	users := a.UserContext.GetAllLockedUsers()

	data := adminPageData{
		Users: users,
	}

	templateRenderError := templateManager.RenderTemplate(w, "AdminPage", data)

	if templateRenderError != nil {
		a.Logger.LogError("Admin", templateRenderError)
	}
}
