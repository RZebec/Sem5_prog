package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
)

/*
	Structure for the Admin handler wrapper.
*/
type AdminWrapper struct {
	Next        HttpHandler
	Config      config.Configuration
	UserContext user.UserContext
	Logger      logging.Logger
}

/*
	The Admin handler wrapper.
*/
func (h AdminWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(shared.AccessTokenCookieName)

	if err != nil {
		h.Logger.LogError("Admin", err)
		return
	}

	isSessionValid, _, _, role, err := h.UserContext.SessionIsValid(cookie.Value)

	if err != nil {
		h.Logger.LogError("Admin", err)
		return
	}

	userIsAdmin := isSessionValid && role == user.Admin

	if userIsAdmin {
		h.Next.ServeHTTP(w, r)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusAccepted)
		return
	}
}