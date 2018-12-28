package api

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core"
	"net/http"
)

type ApiKeyAuthenticationHandler struct {
	Next           core.HttpHandler
	ApiKeyResolver func() string
}

func (h *ApiKeyAuthenticationHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, cookie := range req.Cookies() {
		if cookie.Name == shared.AuthenticationCookieName {
			apiKey := cookie.Value
			if h.ApiKeyResolver() == apiKey {
				h.Next.ServeHTTP(w, req)
			} else {
				w.WriteHeader(401)
			}
			break
		}
	}
	w.WriteHeader(401)
}
