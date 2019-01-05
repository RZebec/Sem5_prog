package api

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core"
	"net/http"
	"strings"
)

type ApiKeyAuthenticationHandler struct {
	Next           core.HttpHandler
	ApiKeyResolver func() string
	AllowedMethod  string
	Logger         logging.Logger
}

func (h *ApiKeyAuthenticationHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != h.AllowedMethod {
		h.Logger.LogInfo("ApiKeyAuthenticationHandler", "Invalid HTTP Method -> 401")
		w.WriteHeader(401)
		return
	}

	cookies := req.Cookies()
	for _, cookie := range cookies {
		if strings.ToLower(cookie.Name) == strings.ToLower(shared.AuthenticationCookieName) {
			apiKey := cookie.Value
			currentApiKey := h.ApiKeyResolver()
			if currentApiKey == apiKey {
				h.Next.ServeHTTP(w, req)
				return
			} else {
				w.WriteHeader(401)
				h.Logger.LogInfo("ApiKeyAuthenticationHandler", "Wrong Api Key -> 401")
				return
			}
		}
	}
	w.WriteHeader(401)
	h.Logger.LogInfo("ApiKeyAuthenticationHandler", "Api key not set -> 401")
}
