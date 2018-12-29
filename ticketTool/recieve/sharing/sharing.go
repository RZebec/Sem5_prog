package sharing

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
)
func SharingAllMails(mails []mail.Mail)[]mail.Acknowledgment{
	acknowledgeIds := make([]mail.Acknowledgment,len(mails))
	for i:=0;i<len(mails);i++{
		fmt.Println("From: "+mails[i].Sender)
		fmt.Println("Send to: "+mails[i].Receiver)
		fmt.Println("Subject: "+mails[i].Subject)
		fmt.Println("Content: "+mails[i].Content)
		acknowledgeIds[i].Id = mails[i].Id
		acknowledgeIds[i].Subject = mails[i].Subject
	}
	return acknowledgeIds
}

func SharingSingleMails(mails []mail.Mail)mail.Acknowledgment{
	 
}

func GetAcknowledges(mails []mail.Mail)[]mail.Acknowledgment{
	acknowledgeIds := make([]mail.Acknowledgment,len(mails))
	for i:=0;i<len(mails);i++{
		acknowledgeIds[i].Id = mails[i].Id
		acknowledgeIds[i].Subject = mails[i].Subject
	}
	return acknowledgeIds
}

