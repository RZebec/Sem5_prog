package userSettings

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"html"
	"net/http"
	"strings"
)

/*
	Structure for the User Settings Page handler.
*/
type SettingsPageHandler struct {
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
	UserContext     userData.UserContext
}

/*
	Structure for the User Settings Page Data.
*/
type userSettingsPageData struct {
	pages.BasePageData
	IsChangeFailed   string
	UserIsOnVacation bool
}

/*
	The User Settings Page handler.
*/
func (u SettingsPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
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

		userId := wrappers.GetUserId(r.Context())

		exist, loggedInUser := u.UserContext.GetUserById(userId)

		userIsOnVacation := false

		if exist {
			userIsOnVacation = loggedInUser.State == userData.OnVacation
		}

		data := userSettingsPageData{
			IsChangeFailed:   isChangeFailed,
			UserIsOnVacation: userIsOnVacation,
		}

		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.UserIsAdmin = wrappers.IsAdmin(r.Context())
		data.Active = "settings"

		templateRenderError := u.TemplateManager.RenderTemplate(w, "UserSettingsPage", data)

		if templateRenderError != nil {
			u.Logger.LogError("SettingsPageHandler", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}

}
