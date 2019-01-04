package index

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
)

/*
	Structure for the Index Page Handler.
*/
type IndexPageHandler struct {
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the Index Page Data.
*/
type indexPageData struct {
	pages.BasePageData
}

/*
	The Index Page handler.
*/
func (i IndexPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		pageData := new(indexPageData)
		pageData.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		pageData.UserIsAdmin = wrappers.IsAdmin(r.Context())
		pageData.Active = "index"

		err := i.TemplateManager.RenderTemplate(w, "IndexPage", pageData)

		if err != nil {
			i.Logger.LogError("Index", err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
