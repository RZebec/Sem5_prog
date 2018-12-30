package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
)

/*
	Structure for the Admin handler wrapper.
*/
type AdminWrapper struct {
	Next        HttpHandler
	Config      config.Configuration
	UserContext user.UserContext
	Logger      logging.Logger
}

/*
	The Admin handler wrapper.
*/
func (h AdminWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIsAdmin := isAdmin(r.Context())

	if userIsAdmin {
		h.Next.ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
	}
}