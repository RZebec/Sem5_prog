package admin

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
	"strconv"
	"strings"
)

/*
	Structure for the Login handler.
*/
type AdminUnlockUserHandler struct {
	UserContext user.UserContext
	Config      config.Configuration
	Logger      logging.Logger
}

/*
	The Unlock user handler.
*/
func (a AdminUnlockUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		URLPath := strings.Split(r.URL.Path, "/")

		userId, idConversionError := strconv.Atoi(URLPath[2])

		if idConversionError != nil {
			a.Logger.LogError("Admin", idConversionError)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		accessTokenCookie, err := r.Cookie(shared.AccessTokenCookieName)

		if err != nil {
			a.Logger.LogError("Login", err)
			return
		}

		unlocked, err := a.UserContext.UnlockAccount(accessTokenCookie.Value, userId)

		if err != nil {
			a.Logger.LogError("Admin", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		if unlocked {
			http.Redirect(w, r, "/admin", http.StatusCreated)
			return
		}
	}
}

