package admin

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
	"strconv"
)

/*
	Structure for the Login handler.
*/
type AdminUnlockUserHandler struct {
	UserContext user.UserContext
	Logger      logging.Logger
}

/*
	The Unlock user handler.
*/
func (a AdminUnlockUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		formId := r.FormValue("userId")

		userId, idConversionError := strconv.Atoi(formId)

		if idConversionError != nil {
			a.Logger.LogError("Admin", idConversionError)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		accessTokenCookie, err := r.Cookie(shared.AccessTokenCookieName)

		if err != nil {
			a.Logger.LogError("Admin", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		unlocked, err := a.UserContext.UnlockAccount(accessTokenCookie.Value, userId)

		if err != nil {
			a.Logger.LogError("Admin", err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}

		if unlocked {
			http.Redirect(w, r, "/admin", http.StatusOK)
		}
	}
}

