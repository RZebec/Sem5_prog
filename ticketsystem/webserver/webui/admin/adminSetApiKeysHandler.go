package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"html"
	"net/http"
	"strings"
)

/*
	Structure for the Login handler.
*/
type AdminSetApiKeysHandler struct {
	ApiConfiguration	config.ApiContext
	Logger      		logging.Logger
}

/*
	The Api key post handler.
*/
func (a AdminSetApiKeysHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		incomingMailApiKey := r.FormValue("incomingMailApiKey")
		outgoingMailApiKey := r.FormValue("outgoingMailApiKey")

		incomingMailApiKey = html.EscapeString(incomingMailApiKey)
		outgoingMailApiKey = html.EscapeString(outgoingMailApiKey)

		if len(incomingMailApiKey) >= 128 && len(outgoingMailApiKey) >= 128 {
			a.ApiConfiguration.ChangeIncomingMailApiKey(incomingMailApiKey)
			a.ApiConfiguration.ChangeOutgoingMailApiKey(outgoingMailApiKey)
		}
	}
}

