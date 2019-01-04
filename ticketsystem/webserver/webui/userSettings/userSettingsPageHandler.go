package userSettings

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
	"strconv"
	"strings"
)

/*
	Structure for the User Settings Page handler.
*/
type UserSettingsPageHandler struct {
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the User Settings Page Data.
*/
type userSettingsPageData struct {
	pages.BasePageData
	IsChangeFailed bool
}

/*
	The User Settings Page handler.
*/
func (u UserSettingsPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		queryValues := r.URL.Query()
		queryValue := queryValues.Get("IsChangeFailed")
		isChangeFailed, parseError := strconv.ParseBool(queryValue)

		if parseError != nil && queryValue != ""{
			u.Logger.LogError("UserSettingsPageHandler", parseError)
			isChangeFailed = false
		}

		data := userSettingsPageData{
			IsChangeFailed: isChangeFailed,
		}

		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.Active = "settings"

		templateRenderError := u.TemplateManager.RenderTemplate(w, "UserSettingsPage", data)

		if templateRenderError != nil {
			u.Logger.LogError("UserSettingsPageHandler", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}

}