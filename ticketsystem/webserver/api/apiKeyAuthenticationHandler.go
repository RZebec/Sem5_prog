package api

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core"
	"net/http"
	"strings"
)

type ApiKeyAuthenticationHandler struct {
	Next           core.HttpHandler
	ApiKeyResolver func() string
	AllowedMethod string
}

func (h *ApiKeyAuthenticationHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != h.AllowedMethod {
		w.WriteHeader(401)
		return
	}

	cookies := req.Cookies()
	for _, cookie := range cookies{
		if strings.ToLower(cookie.Name) == strings.ToLower(shared.AuthenticationCookieName) {
			apiKey := cookie.Value
			currentApiKey := h.ApiKeyResolver()
			if currentApiKey == apiKey {
				h.Next.ServeHTTP(w, req)
				return
			} else {
				w.WriteHeader(401)
				return
			}
			break
		}
	}
	w.WriteHeader(401)
}
