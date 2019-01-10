// 5894619, 6720876, 9793350
package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"html"
	"net/http"
)

/*
	Structure for the Admin handler.
*/
type PageHandler struct {
	UserContext           userData.UserContext
	Logger                logging.Logger
	TemplateManager       templateManager.TemplateContext
	GetIncomingMailApiKey func() string
	GetOutgoingMailApiKey func() string
}

/*
	Structure for the Admin Page Data.
*/
type adminPageData struct {
	Users              []userData.User
	IncomingMailApiKey string
	OutgoingMailApiKey string
	IsChangeFailed     string
	pages.BasePageData
}

/*
	The Admin Page handler.
*/
func (a PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		queryValues := r.URL.Query()
		queryValue := queryValues.Get("IsChangeFailed")

		queryValue = html.EscapeString(queryValue)

		isChangeFailed := queryValue

		if queryValue == "yes" || queryValue == "no" {
			isChangeFailed = queryValue
		} else {
			isChangeFailed = "NotSet"
		}

		users := a.UserContext.GetAllLockedUsers()

		incomingMailApiKey := a.GetIncomingMailApiKey()
		outgoingMailApiKey := a.GetOutgoingMailApiKey()

		data := adminPageData{
			Users:              users,
			IncomingMailApiKey: incomingMailApiKey,
			OutgoingMailApiKey: outgoingMailApiKey,
			IsChangeFailed:     isChangeFailed,
		}
		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.Active = "admin"

		templateRenderError := a.TemplateManager.RenderTemplate(w, "AdminPage", data)

		if templateRenderError != nil {
			a.Logger.LogError("AdminPageHandler", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
