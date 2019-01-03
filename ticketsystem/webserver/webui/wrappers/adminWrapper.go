package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
)

/*
	Structure for the Admin handler wrapper.
*/
type AdminWrapper struct {
	Next        HttpHandler
	UserContext user.UserContext
	Logger      logging.Logger
}

/*
	The Admin handler wrapper.
*/
func (h AdminWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIsAdmin := IsAdmin(r.Context())

	if userIsAdmin {
		h.Next.ServeHTTP(w, r)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}
}
