// 5894619, 6720876, 9793350
package core

import (
	"net/http"
)

type HttpHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
