package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"encoding/json"
	"fmt"
	"net/http"
)

type IncomingMailHandler struct {
}

func (h *IncomingMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t []mail.Mail
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t)
	}
	w.WriteHeader(200)
}
