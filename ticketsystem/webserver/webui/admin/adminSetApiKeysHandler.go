package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"html"
	"net/http"
)

/*
	Structure for the Login handler.
*/
type SetApiKeysHandler struct {
	Logger           logging.Logger
	ChangeIncomingMailApiKey func(newKey string) error
	ChangeOutgoingMailApiKey func(newKey string) error
}

/*
	The Api key post handler.
*/
func (a SetApiKeysHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		incomingMailApiKey := r.FormValue("incomingMailApiKey")
		outgoingMailApiKey := r.FormValue("outgoingMailApiKey")

		incomingMailApiKey = html.EscapeString(incomingMailApiKey)
		outgoingMailApiKey = html.EscapeString(outgoingMailApiKey)

		if len(incomingMailApiKey) >= 128 && len(outgoingMailApiKey) >= 128 {
			err := a.ChangeIncomingMailApiKey(incomingMailApiKey)
			if err != nil {
				http.Redirect(w, r, "/admin?IsChangeFailed=yes", http.StatusInternalServerError)
				a.Logger.LogError("SetApiKeysHandler", err)
				return
			}
			err = a.ChangeOutgoingMailApiKey(outgoingMailApiKey)
			if err != nil {
				http.Redirect(w, r, "/admin?IsChangeFailed=yes", http.StatusInternalServerError)
				a.Logger.LogError("SetApiKeysHandler", err)
				return
			}
			a.Logger.LogInfo("SetApiKeysHandler", "Api Keys changed")
			http.Redirect(w, r, "/admin?IsChangeFailed=no", http.StatusFound)
		} else {
			http.Redirect(w, r, "/admin?IsChangeFailed=yes", http.StatusBadRequest)
		}
	}
}
