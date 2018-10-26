package webui

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"net/http"
	"strings"
)

type LogoutHandler struct {
	UserManager session.UserManager
}

func (l LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		c, _ := r.Cookie("Access-Token")

		cookie := Cookie{Name: "Access-Token", Value: c.Value}

		if cookie.Value != "" {
			l.UserManager.Logout(cookie.Value)
			cookie.RemoveCookie(w, r)
			http.Redirect(w, r, "/", 302)
		}
	}
}
