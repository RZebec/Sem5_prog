// 5894619, 6720876, 9793350
package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"strings"
)

// Inspired by: https://stackoverflow.com/questions/1027395/detecting-outlook-autoreply-out-of-office-emails

type AutomaticRepliesFilter interface {
	IsAutomaticResponse(mail mailData.Mail) bool
}

/*
	A mail filter struct.
*/
type RepliesFilter struct {
}

/*
	Check if a mail is a automatic response.
*/
func (r *RepliesFilter) IsAutomaticResponse(mail mailData.Mail) bool {
	for _, header := range mail.Headers {
		if strings.Contains(strings.ToLower(header), "x-autorespond") {
			return true
		}
		if strings.Contains(strings.ToLower(header), "precedence") {
			splitted := strings.Split(header, ":")
			if len(splitted) == 2 {
				if strings.ToLower(splitted[1]) == "auto_reply" || strings.ToLower(splitted[1]) == "bulk" || strings.ToLower(splitted[1]) == "junk" {
					return true
				}
			}
		}
		if strings.Contains(strings.ToLower(header), "auto-submitted") {
			splitted := strings.Split(header, ":")
			if len(splitted) == 2 {
				if strings.ToLower(splitted[1]) == "auto-replied" {
					return true
				}
			}
		}
	}

	return false
}
