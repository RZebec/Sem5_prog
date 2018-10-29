package login

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
	"strings"
)

/*
	Structure for the Login handler.
*/
type LoginHandler struct {
	UserContext user.UserContext
	Config      config.Configuration
}

/*
	The Login handler.
*/
func (l LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		userName := r.FormValue("userName")
		password := r.FormValue("password")

		success, token, err := l.UserContext.Login(userName, password)

		if err != nil {
			// TODO: handle error
		}

		if success {
			helpers.SetCookie(w, r, l.Config.AccessTokenCookieName, token)
			http.Redirect(w, r, "/", 302)
		} else {
			http.Redirect(w, r, "/login?IsLoginFailed=true", 302)
		}
	}
}
