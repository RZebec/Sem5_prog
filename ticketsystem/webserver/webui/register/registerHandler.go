package register

import (
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
	Structure for the Register handler.
*/
type RegisterHandler struct {
	UserContext user.UserContext
	Config      config.Configuration
	Logger      logging.Logger
}

/*
	Structure for the Register Page Data.
*/
type registerPageData struct {
	IsRegisteringFailed bool
}

/*
	The Login handler.
*/
func (l RegisterHandler) ServeHTTPPostRegisteringData(w http.ResponseWriter, r *http.Request) {
	// TODO: Verification step for the user needs to be implemented here
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		userName := r.FormValue("userName")
		password := r.FormValue("password")

		success, err := l.UserContext.Register(userName, password, firstName, lastName)

		if err != nil {
			l.Logger.LogError("Register", err)
		}

		if success {
			http.Redirect(w, r, "/login", 302)
		} else {
			http.Redirect(w, r, "/register?IsRegisteringFailed=true", 302)
		}
	}
}

/*
	The Register Page handler.
*/
func (l RegisterHandler) ServeHTTPGetRegisterPage(w http.ResponseWriter, r *http.Request) {

	// Checks if the User is already logged in and if so redirects him to the start page
	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, l.UserContext, l.Config.AccessTokenCookieName, l.Logger)

	if isUserLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	queryValues := r.URL.Query()
	isRegisteringFailed, err := strconv.ParseBool(queryValues.Get("IsRegisteringFailed"))

	if err != nil {
		l.Logger.LogError("Register", err)
	}

	data := registerPageData{
		IsRegisteringFailed: isRegisteringFailed,
	}

	err = templateManager.RenderTemplate(w, "RegisterPage", data)

	if err != nil {
		l.Logger.LogError("Register", err)
	}
}
