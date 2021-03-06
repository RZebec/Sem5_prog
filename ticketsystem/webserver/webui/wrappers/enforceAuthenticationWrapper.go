// 5894619, 6720876, 9793350
package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
)

/*
	Structure for the authentication handler wrapper.
*/
type EnforceAuthenticationWrapper struct {
	Next        HttpHandler
	Config      config.WebServerConfiguration
	UserContext userData.UserContext
	Logger      logging.Logger
}

/*
	The Authentication handler wrapper.
*/
func (h EnforceAuthenticationWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	userIsLoggedIn, isAdmin, token, userId := UserIsLoggedInCheck(r, h.UserContext, shared.AccessTokenCookieName, h.Logger)

	if userIsLoggedIn {
		newToken, err := h.UserContext.RefreshToken(token)

		if err != nil {
			h.Logger.LogError("EnforceAuthenticationWrapper", err)
		}

		helpers.SetCookie(w, shared.AccessTokenCookieName, newToken)

		ctx := NewContextWithAuthenticationInfo(r.Context(), userIsLoggedIn, isAdmin, userId, newToken)
		h.Next.ServeHTTP(w, r.WithContext(ctx))
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}
