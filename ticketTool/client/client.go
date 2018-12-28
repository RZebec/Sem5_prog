package client

import "de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data"

type Client interface {
	SendMails(mails []data.Mail) error
	ReceiveMails() ([]data.Mail, error)
	AcknowledgeMails(mailIds []int) error
}
