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
	Next HttpHandler
	AccessTokenCookie helpers.Cookie
	SessionManager session.SessionManager
}

func (h *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	UserIsLoggedIn := false

	c, err := r.Cookie(h.AccessTokenCookie.Name)

	if err == nil {
		UserIsLoggedIn, _, _, err = h.SessionManager.SessionIsValid(c.Value)
	} else {
		// TODO: Handle the Error
	}

	if UserIsLoggedIn {
		h.Next.ServeHTTP(w, r)
	} else {
		accessDeniedHandler := accessdenied.AccessDeniedPageHandler{}
		accessDeniedHandler.ServeHTTP(w, r)
	}
}