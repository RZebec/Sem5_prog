package wrappers

import "net/http"

/*
	Interface for the Http Handler.
*/
type HttpHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}