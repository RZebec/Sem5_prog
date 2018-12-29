package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
)

/*
	Structure for the authentication handler wrapper.
*/
type AuthenticationWrapper struct {
	Next        HttpHandler
	Config      config.Configuration
	UserContext user.UserContext
	Logger      logging.Logger
}

/*
	The Authentication handler wrapper.
*/
func (h AuthenticationWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	userIsLoggedIn, token := helpers.UserIsLoggedInCheck(r, h.UserContext, shared.AccessTokenCookieName, h.Logger)

	if userIsLoggedIn {
		newToken, err := h.UserContext.RefreshToken(token)

		if err != nil {
			h.Logger.LogError("Login", err)
		}

		helpers.SetCookie(w, r, shared.AccessTokenCookieName, newToken)

		h.Next.ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
	}
}
