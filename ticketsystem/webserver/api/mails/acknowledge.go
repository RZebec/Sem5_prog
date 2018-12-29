package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
	A handler for acknowledgment of mails.
*/
type AcknowledgeMailHandler struct {
	Logger logging.Logger
}

/*
	Handling the acknowledgement of mails.
*/
func (h *AcknowledgeMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t []mail.Acknowledgment
	err := decoder.Decode(&t)
	if err != nil {
		h.Logger.LogError("AcknowledgeMailHandler", err)
		w.WriteHeader(500)
	} else {
		h.Logger.LogInfo("AcknowledgeMailHandler", "Number of acknowledged mails: "+strconv.Itoa(len(t)))
	}
	w.WriteHeader(200)
}
