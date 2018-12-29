package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
	A handler for acknowledgment of mails.
 */
type AcknowledgeMailHandler struct {
}

/*
	Handling the acknowledgement of mails.
 */
func (h *AcknowledgeMailHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t []mail.Acknowledgment
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t)
	}
	w.WriteHeader(200)
}