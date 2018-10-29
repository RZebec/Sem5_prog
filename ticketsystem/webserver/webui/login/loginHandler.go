package login

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"fmt"
	"net/http"
	"strings"
)

type LoginHandler struct {
	UserManager session.UserManager
	Config 	config.Configuration
}

func (l LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		userName := r.FormValue("userName")
		password := r.FormValue("password")

		success, token, err := l.UserManager.Login(userName, password)

		if err != nil {
			fmt.Println(err.Error())
		}

		if success {
			helpers.SetCookie(w, r, l.Config.AccessTokenCookieName, token)
			http.Redirect(w, r, "/", 302)
		} else {
			http.Redirect(w, r, "/login?IsLoginFailed=true", 302)
		}
	}
}
