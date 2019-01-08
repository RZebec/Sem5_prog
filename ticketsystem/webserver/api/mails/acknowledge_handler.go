// 5894619, 6720876, 9793350
package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
	A handler for acknowledgment of mails.
*/
type AcknowledgeMailHandler struct {
	Logger      logging.Logger
	MailContext mailData.MailContext
}

/*
	Handling the acknowledgement of mails.
*/
func (h *AcknowledgeMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data []mailData.Acknowledgment
	err := decoder.Decode(&data)
	if err != nil {
		h.Logger.LogError("AcknowledgeMailHandler", err)
		w.WriteHeader(400)
	} else {
		err := h.MailContext.AcknowledgeMails(data)
		if err != nil {
			h.Logger.LogError("AcknowledgeMailHandler", err)
			w.WriteHeader(500)
			return
		}
		h.Logger.LogInfo("AcknowledgeMailHandler", "Number of acknowledged mails: "+strconv.Itoa(len(data)))
	}
	w.WriteHeader(200)
}
