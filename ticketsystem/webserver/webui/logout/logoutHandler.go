package logout

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
	"strings"
)

/*
	Structure for the Logout handler.
*/
type UserLogoutHandler struct {
	UserContext userData.UserContext
	Logger      logging.Logger
}

/*
	The Logout handler.
*/
func (l UserLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		cookie, err := r.Cookie(shared.AccessTokenCookieName)

		if err != nil {
			l.Logger.LogError("UserLogoutHandler", err)
			helpers.RemoveCookie(w, shared.AccessTokenCookieName)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		token := cookie.Value
		l.UserContext.Logout(token)
		helpers.RemoveCookie(w, shared.AccessTokenCookieName)
		http.Redirect(w, r, "/", http.StatusFound)
		l.Logger.LogInfo("UserLogoutHandler","User logged out" )

		return
	}
}
