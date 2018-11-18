package login

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"net/http"
	"strconv"
)

/*
	Structure for the Login Page Handler.
*/
type LoginPageHandler struct {
	Config      config.Configuration
	UserContext user.UserContext
}

/*
	Structure for the Login Page Data.
*/
type loginPageData struct {
	IsLoginFailed bool
}

/*
	The Login Page handler.
*/
func (l LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Checks if the User is already logged in and if so redirects him to the start page
	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, l.UserContext, l.Config.AccessTokenCookieName)

	if isUserLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	queryValues := r.URL.Query()
	isLoginFailed, parseError := strconv.ParseBool(queryValues.Get("IsLoginFailed"))

	if parseError != nil {
		// TODO: Handle error
	}

	data := loginPageData{
		IsLoginFailed: isLoginFailed,
	}

	templateRenderError := templateManager.RenderTemplate(w, "LoginPage", data)

	if templateRenderError != nil {
		// TODO: Handle error
	}
}
