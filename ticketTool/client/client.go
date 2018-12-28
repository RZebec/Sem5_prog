package client

import "de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"

type Client interface {
	SendMails(mails []mail.Mail) error
	ReceiveMails() ([]mail.Mail, error)
	AcknowledgeMails(mailIds []int) error
}
