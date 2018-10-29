package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/accessdenied"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"net/http"
)

type HttpHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type AuthenticationHandler struct {
	Next              HttpHandler
	AccessTokenCookie helpers.Cookie
	SessionManager    session.SessionManager
}

func (h *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	userIsLoggedIn, token := helpers.UserIsLoggedInCheck(r, h.SessionManager, h.AccessTokenCookie)

	if userIsLoggedIn {
		newToken, err := h.SessionManager.RefreshToken(token)

		if err != nil {
			panic(err)
		}

		h.AccessTokenCookie.Value = newToken
		h.AccessTokenCookie.SetCookie(w, r)

		h.Next.ServeHTTP(w, r)
	} else {
		accessDeniedHandler := accessdenied.AccessDeniedPageHandler{}
		accessDeniedHandler.ServeHTTP(w, r)
	}
}
