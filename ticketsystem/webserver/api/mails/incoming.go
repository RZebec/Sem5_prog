package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
	A handler for incoming mails.
*/
type IncomingMailHandler struct {
	Logger logging.Logger
}

/*
	Handling the incoming mails.
*/
func (h *IncomingMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t []mail.Mail
	err := decoder.Decode(&t)
	if err != nil {
		h.Logger.LogError("IncomingMailHandler", err)
		w.WriteHeader(500)
	} else {
		h.Logger.LogInfo("IncomingMailHandler", "Number of received mails: "+strconv.Itoa(len(t)))
	}
	w.WriteHeader(200)
}
