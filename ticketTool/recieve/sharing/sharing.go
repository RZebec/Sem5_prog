package sharing

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
)

func SharingAllMails(mails []mail.Mail) []mail.Acknowledgment {
	acknowledge := make([]mail.Acknowledgment, len(mails))
	for i := 0; i < len(mails); i++ {
		fmt.Println("From: " + mails[i].Sender)
		fmt.Println("Send to: " + mails[i].Receiver)
		fmt.Println("Subject: " + mails[i].Subject)
		fmt.Println("Content: " + mails[i].Content)
		acknowledge[i].Id = mails[i].Id
		acknowledge[i].Subject = mails[i].Subject
	}
	return acknowledge
}

func SharingSingleMails(mails []mail.Mail) []mail.Acknowledgment {
	acknowledge := make([]mail.Acknowledgment, 1)
	fmt.Println("Specify Mail by Subject: ")
	answer := inputOutput.ReadEntry()
	for i := 0; i < len(mails); i++ {
		if answer == mails[i].Subject {
			fmt.Println("From: " + mails[i].Sender)
			fmt.Println("Send to: " + mails[i].Receiver)
			fmt.Println("Subject: " + mails[i].Subject)
			fmt.Println("Content: " + mails[i].Content)
			acknowledge[0].Id = mails[i].Id
			acknowledge[0].Subject = mails[i].Subject
			return acknowledge
		}
	}
	fmt.Println("Subject not found")
	return acknowledge
}

func GetAcknowledges(mails []mail.Mail) []mail.Acknowledgment {
	acknowledge := make([]mail.Acknowledgment, len(mails))
	for i := 0; i < len(mails); i++ {
		acknowledge[i].Id = mails[i].Id
		acknowledge[i].Subject = mails[i].Subject
	}
	return acknowledge
}
