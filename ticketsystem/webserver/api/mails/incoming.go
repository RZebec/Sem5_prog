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
	MailContext mail.MailContext
}

/*
	Handling the incoming mails.
*/
func (h *IncomingMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data []mail.Mail
	err := decoder.Decode(&data)
	if err != nil {
		h.Logger.LogError("IncomingMailHandler", err)
		w.WriteHeader(500)
	} else {
		h.Logger.LogInfo("IncomingMailHandler", "Number of received mails: "+strconv.Itoa(len(data)))
	}
	w.WriteHeader(200)
}
