package login

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"html"
	"net/http"
	"strings"
)

/*
	Structure for the User Login handler.
*/
type UserLoginHandler struct {
	UserContext     userData.UserContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	The User Login handler.
*/
func (l UserLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		userName := r.FormValue("userName")
		password := r.FormValue("password")

		userName = html.EscapeString(userName)
		password = html.EscapeString(password)

		success, token, err := l.UserContext.Login(userName, password)

		if err != nil {
			l.Logger.LogError("UserLoginHandler", err)
		}

		if success {
			helpers.SetCookie(w, shared.AccessTokenCookieName, token)
			http.Redirect(w, r, "/", 302)
			l.Logger.LogInfo("UserLoginHandler", "User logged in")
		} else {
			http.Redirect(w, r, "/login?IsLoginFailed=true", 302)
			l.Logger.LogInfo("UserLoginHandler", "Login failed")
		}
	}
}
