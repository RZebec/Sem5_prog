package logout

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
	"strings"
)

type LogoutHandler struct {
	UserManager       session.UserManager
	AccessTokenCookie helpers.Cookie
}

func (l LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		l.UserManager.Logout(l.AccessTokenCookie.Value)
		l.AccessTokenCookie.RemoveCookie(w, r)
		http.Redirect(w, r, "/", 302)
	}
}
