package index

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"net/http"
)

/*
	Structure for the Index Page Handler.
*/
type IndexPageHandler struct {
	Config      config.Configuration
	UserContext user.UserContext
	Logger		logging.Logger
}

/*
	Structure for the Index Page Data.
*/
type indexPageData struct {
	IsUserLoggedIn bool
}

/*
	The Index Page handler.
*/
func (i IndexPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templateManager.LoadTemplates(i.Logger)

	err := templateManager.RenderTemplate(w, "IndexPage", nil)

	if err != nil {
		i.Logger.LogError("Index", err)
	}
}
