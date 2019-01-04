package userSettings

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"html"
	"net/http"
	"strconv"
	"strings"
)

/*
	Structure for the Toggle Vacation Mode handler.
*/
type ToggleVacationModeHandler struct {
	UserContext     user.UserContext
	Logger          logging.Logger
}

/*
	The Toggle Vacation Mode handler.
*/
func (v ToggleVacationModeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		vacationMode := r.FormValue("vacationMode")

		vacationMode = html.EscapeString(vacationMode)

		vacation, parseErr := strconv.ParseBool(vacationMode)

		if parseErr != nil {
			v.Logger.LogError("ToggleVacationModeHandler", parseErr)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		accessTokenCookie, err := r.Cookie(shared.AccessTokenCookieName)

		if err != nil {
			v.Logger.LogError("ToggleVacationModeHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		if vacation {
			err := v.UserContext.EnableVacationMode(accessTokenCookie.Value)
			if err != nil {
				v.Logger.LogError("ToggleVacationModeHandler", err)
				http.Redirect(w, r, "/", http.StatusBadRequest)
				return
			}
		} else {
			err := v.UserContext.DisableVacationMode(accessTokenCookie.Value)
			if err != nil {
				v.Logger.LogError("ToggleVacationModeHandler", err)
				http.Redirect(w, r, "/", http.StatusBadRequest)
				return
			}
		}

		http.Redirect(w, r, "/user_settings", 302)
	}
}
