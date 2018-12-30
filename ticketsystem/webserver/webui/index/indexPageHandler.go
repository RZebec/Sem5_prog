package index

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"net/http"
)

/*
	Structure for the Index Page Handler.
*/
type IndexPageHandler struct {
	Logger		logging.Logger
	TemplateManager	templateManager.ITemplateManager
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
	err := i.TemplateManager.RenderTemplate(w, "IndexPage", nil)

	if err != nil {
		i.Logger.LogError("Index", err)
	}
}
