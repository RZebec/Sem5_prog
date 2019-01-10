// 5894619, 6720876, 9793350
package login

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
	Structure for the Login Page handler.
*/
type PageHandler struct {
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the Login Page Data.
*/
type loginPageData struct {
	pages.BasePageData
	IsLoginFailed bool
}

/*
	The Login Page handler.
*/
func (l PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		// Checks if the User is already logged in and if so redirects him to the start page
		isUserLoggedIn := wrappers.IsAuthenticated(r.Context())
		userIsAdmin := wrappers.IsAdmin(r.Context())

		if isUserLoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		queryValues := r.URL.Query()
		queryValue := queryValues.Get("IsLoginFailed")
		isLoginFailed, parseError := strconv.ParseBool(queryValue)

		if parseError != nil && queryValue != "" {
			l.Logger.LogError("PageHandler", parseError)
		}

		data := loginPageData{
			IsLoginFailed: isLoginFailed,
		}
		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.UserIsAdmin = userIsAdmin
		data.Active = "login"

		templateRenderError := l.TemplateManager.RenderTemplate(w, "LoginPage", data)

		if templateRenderError != nil {
			l.Logger.LogError("PageHandler", templateRenderError)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
