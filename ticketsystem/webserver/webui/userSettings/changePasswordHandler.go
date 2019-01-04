package userSettings

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"html"
	"net/http"
	"strings"
)

/*
	Structure for the Login handler.
*/
type ChangePasswordHandler struct {
	UserContext     user.UserContext
	Logger          logging.Logger
}

/*
	The Login handler.
*/
func (c ChangePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")

		oldPassword = html.EscapeString(oldPassword)
		newPassword = html.EscapeString(newPassword)

		accessTokenCookie, err := r.Cookie(shared.AccessTokenCookieName)

		if err != nil {
			c.Logger.LogError("ChangePasswordHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		isChanged, err := c.UserContext.ChangePassword(accessTokenCookie.Value, oldPassword, newPassword)

		if err != nil {
			c.Logger.LogError("ChangePasswordHandler", err)
			http.Redirect(w, r, "/user_settings?IsChangeFailed=true", 302)
		}

		if isChanged {
			http.Redirect(w, r, "/user_settings", 302)
		} else {
			http.Redirect(w, r, "/user_settings?IsChangeFailed=true", 302)
		}
	}
}
