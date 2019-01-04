package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
	A handler for outgoing mails.
*/
type OutgoingMailHandler struct {
	Logger      logging.Logger
	MailContext mail.MailContext
}

/*
	Handling the outgoing mails.
*/
func (h *OutgoingMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mails, err := h.MailContext.GetUnsentMails()
	if err != nil {
		w.WriteHeader(500)
		h.Logger.LogError("OutgoingMailHandler", err)
		return
	}
	jsonData, err := json.Marshal(mails)
	if err != nil {
		w.WriteHeader(500)
		h.Logger.LogError("OutgoingMailHandler", err)
		return
	}
	w.WriteHeader(200)
	w.Write(jsonData)
	h.Logger.LogInfo("OutgoingMailHandler", "Number of outgoing mails: "+strconv.Itoa(len(mails)))
}
