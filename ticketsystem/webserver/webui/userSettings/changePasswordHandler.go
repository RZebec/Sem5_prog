package userSettings

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"html"
	"net/http"
	"strings"
)

/*
	Structure for the Change Password handler.
*/
type ChangePasswordHandler struct {
	UserContext userData.UserContext
	Logger      logging.Logger
}

/*
	The Change Password handler.
*/
func (c ChangePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")

		oldPassword = html.EscapeString(oldPassword)
		newPassword = html.EscapeString(newPassword)

		accessToken := wrappers.GetUserToken(r.Context())

		isChanged, err := c.UserContext.ChangePassword(accessToken, oldPassword, newPassword)

		if err != nil {
			c.Logger.LogError("ChangePasswordHandler", err)
			http.Redirect(w, r, "/user_settings?IsChangeFailed=yes", 302)
		}

		if isChanged {
			http.Redirect(w, r, "/user_settings?IsChangeFailed=no", 302)
		} else {
			http.Redirect(w, r, "/user_settings?IsChangeFailed=yes", 302)
		}
	}
}
