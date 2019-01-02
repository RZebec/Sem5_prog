package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
)

/*
	Structure for the authentication handler wrapper.
*/
type AddAuthenticationInfoWrapper struct {
	Next        HttpHandler
	UserContext user.UserContext
	Logger      logging.Logger
}

/*
	The Authentication handler wrapper.
*/
func (h AddAuthenticationInfoWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	userIsLoggedIn, isAdmin, token := UserIsLoggedInCheck(r, h.UserContext, shared.AccessTokenCookieName, h.Logger)

	if userIsLoggedIn {
		newToken, err := h.UserContext.RefreshToken(token)

		if err != nil {
			h.Logger.LogError("AddAuthenticationInfoWrapper", err)
		}

		helpers.SetCookie(w, shared.AccessTokenCookieName, newToken)
	}
	ctx := newContextWithAuthenticationInfo(r.Context(), userIsLoggedIn, isAdmin)
	h.Next.ServeHTTP(w, r.WithContext(ctx))
}