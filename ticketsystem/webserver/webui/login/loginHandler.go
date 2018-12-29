package login

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"net/http"
	"strconv"
	"strings"
)

/*
	Structure for the Login handler.
*/
type LoginHandler struct {
	UserContext user.UserContext
	Config      config.Configuration
	Logger      logging.Logger
}

/*
	Structure for the Login Page Data.
*/
type loginPageData struct {
	IsLoginFailed bool
}

/*
	The Login handler.
*/
func (l LoginHandler) ServeHTTPPostLoginData(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		userName := r.FormValue("userName")
		password := r.FormValue("password")

		success, token, err := l.UserContext.Login(userName, password)

		if err != nil {
			l.Logger.LogError("Login", err)
		}

		if success {
			helpers.SetCookie(w, r, shared.AccessTokenCookieName, token)
			http.Redirect(w, r, "/", 302)
		} else {
			http.Redirect(w, r, "/login?IsLoginFailed=true", 302)
		}
	}
}

/*
	The Login Page handler.
*/
func (l LoginHandler) ServeHTTPGetLoginPage(w http.ResponseWriter, r *http.Request) {

	// Checks if the User is already logged in and if so redirects him to the start page
	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, l.UserContext, shared.AccessTokenCookieName, l.Logger)

	if isUserLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	queryValues := r.URL.Query()
	isLoginFailed, parseError := strconv.ParseBool(queryValues.Get("IsLoginFailed"))

	if parseError != nil {
		l.Logger.LogError("Login", parseError)
	}

	data := loginPageData{
		IsLoginFailed: isLoginFailed,
	}

	templateRenderError := templateManager.RenderTemplate(w, "LoginPage", data)

	if templateRenderError != nil {
		l.Logger.LogError("Login", templateRenderError)
	}
}
