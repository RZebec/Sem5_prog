package logout

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
	"strings"
)

/*
	Structure for the Logout handler.
*/
type LogoutHandler struct {
	UserContext user.UserContext
	Config      config.Configuration
	Logger      logging.Logger
}

/*
	The Logout handler.
*/
func (l LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		cookie, err := r.Cookie(shared.AccessTokenCookieName)

		if err != nil {
			l.Logger.LogError("Logout", err)
			return
		}

		token := cookie.Value
		l.UserContext.Logout(token)
		helpers.SetCookie(w, shared.AccessTokenCookieName, "")
		http.Redirect(w, r, "/", 302)

		return
	}
}
