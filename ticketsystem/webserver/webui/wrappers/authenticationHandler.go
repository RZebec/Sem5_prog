package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/accessdenied"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
)

type HttpHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type AuthenticationHandler struct {
	Next              HttpHandler
	Config 			config.Configuration
	UserContext user.UserContext
}

// todo: no pointer maybe
func (h AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	userIsLoggedIn, token := helpers.UserIsLoggedInCheck(r, h.UserContext, h.Config.AccessTokenCookieName)

	if userIsLoggedIn {
		newToken, err := h.UserContext.RefreshToken(token)

		if err != nil {
			panic(err)
		}

		helpers.SetCookie(w, r,  h.Config.AccessTokenCookieName, newToken)

		h.Next.ServeHTTP(w, r)
	} else {
		accessDeniedHandler := accessdenied.AccessDeniedPageHandler{}
		accessDeniedHandler.ServeHTTP(w, r)
	}
}
