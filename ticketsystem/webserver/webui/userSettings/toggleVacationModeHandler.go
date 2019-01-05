package userSettings

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"html"
	"net/http"
	"strconv"
	"strings"
)

/*
	Structure for the Toggle Vacation Mode handler.
*/
type ToggleVacationModeHandler struct {
	UserContext userData.UserContext
	Logger      logging.Logger
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

		accessToken := wrappers.GetUserToken(r.Context())

		if vacation {
			err := v.UserContext.EnableVacationMode(accessToken)
			if err != nil {
				v.Logger.LogError("ToggleVacationModeHandler", err)
				http.Redirect(w, r, "/", http.StatusBadRequest)
				return
			}
		} else {
			err := v.UserContext.DisableVacationMode(accessToken)
			if err != nil {
				v.Logger.LogError("ToggleVacationModeHandler", err)
				http.Redirect(w, r, "/", http.StatusBadRequest)
				return
			}
		}

		http.Redirect(w, r, "/user_settings", 302)
	}
}
