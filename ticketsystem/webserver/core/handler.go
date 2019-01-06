package core

import (
	"net/http"
)

type HttpHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
