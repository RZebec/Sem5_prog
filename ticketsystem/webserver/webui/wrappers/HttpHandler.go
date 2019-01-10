// 5894619, 6720876, 9793350
package wrappers

import "net/http"

/*
	Interface for the Http Handler.
*/
type HttpHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
