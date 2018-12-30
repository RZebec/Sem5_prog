package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
	"strings"
)

/*
	Structure for the Login handler.
*/
type AdminSetApiKeysHandler struct {
	UserContext user.UserContext
	Config      config.Configuration
	Logger      logging.Logger
}

/*
	The Api key post handler.
*/
func (a AdminSetApiKeysHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		// incomingMailApiKey := r.FormValue("incomingMailApiKey")
		// outgoingMailApiKey := r.FormValue("outgoingMailApiKey")

		// incomingMailApiKey = html.EscapeString(incomingMailApiKey)
		// outgoingMailApiKey = html.EscapeString(outgoingMailApiKey)


		// TODO: Do something with the keys
	}
}

