package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

/*
	A handler for outgoing mails.
*/
type OutgoingMailHandler struct {
	Logger logging.Logger
}

func getTestEmails() []mail.Mail {
	var eMails []mail.Mail
	eMails = append(eMails, mail.Mail{Id: "testId1", Sender: "test1@test.de", Receiver: "testReceiver1@test.de",
		Subject: "TestSubject1", Content: "testContent1", SentTime: time.Now().Unix()})
	eMails = append(eMails, mail.Mail{Id: "testId2", Sender: "test2@test.de", Receiver: "testReceiver2@test.de",
		Subject: "TestSubject2", Content: "testContent2", SentTime: time.Now().Unix()})
	eMails = append(eMails, mail.Mail{Id: "testId3", Sender: "test3@test.de", Receiver: "testReceiver3@test.de",
		Subject: "TestSubject3", Content: "testContent3", SentTime: time.Now().Unix()})
	return eMails
}

/*
	Handling the outgoing mails.
*/
func (h *OutgoingMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	jsonData, err := json.Marshal(getTestEmails())
	if err != nil {
		w.WriteHeader(500)
		h.Logger.LogError("OutgoingMailHandler", err)
		return
	}
	w.WriteHeader(200)
	w.Write(jsonData)
	h.Logger.LogInfo("OutgoingMailHandler", "Number of outgoing mails: "+strconv.Itoa(len(getTestEmails())))
}
