package logout

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
	"strings"
)

type LogoutHandler struct {
	UserManager       session.UserManager
	Config	config.Configuration
}

func (l LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		cookie, err := r.Cookie(l.Config.AccessTokenCookieName)

		if err == nil {
			token := cookie.Value
			l.UserManager.Logout(token)
			helpers.RemoveCookie(w, r, l.Config.AccessTokenCookieName)
			http.Redirect(w, r, "/", 302)
		}
		// Todo: error handling
	}
}
