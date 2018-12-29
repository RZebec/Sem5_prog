package sharing

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
)

func ShareAllMails(mails []mail.Mail) []mail.Acknowledgment {
	acknowledge := make([]mail.Acknowledgment, len(mails))
	for i := 0; i < len(mails); i++ {
		fmt.Println("From: " + mails[i].Sender)
		fmt.Println("Send to: " + mails[i].Receiver)
		fmt.Println("Subject: " + mails[i].Subject)
		fmt.Println("Content: " + mails[i].Content)
		fmt.Println("is sended")
		acknowledge[i].Id = mails[i].Id
		acknowledge[i].Subject = mails[i].Subject
	}
	return acknowledge
}

func ShareSingleMails(mails *[]mail.Mail) []mail.Acknowledgment {
	showAllMails(*mails)
	acknowledge := make([]mail.Acknowledgment, 1)
	fmt.Println("Specify Mail by Subject: ")
	answer := inputOutput.ReadEntry()
	for i := 0; i < len(*mails); i++ {
		if answer == (*mails)[i].Subject {
			fmt.Println("From: " + (*mails)[i].Sender)
			fmt.Println("Send to: " + (*mails)[i].Receiver)
			fmt.Println("Subject: " + (*mails)[i].Subject)
			fmt.Println("Content: " + (*mails)[i].Content)
			fmt.Println("is sended")
			acknowledge[0].Id = (*mails)[i].Id
			acknowledge[0].Subject = (*mails)[i].Subject
			newMails := DeleteFromArray(*mails, (*mails)[i])
			mails = &newMails
			return acknowledge
		}
	}
	fmt.Println("Subject not found")
	return acknowledge
}

func showAllMails(mails []mail.Mail) {
	fmt.Println("All Mails: ")
	for i := 0; i < len(mails); i++ {
		fmt.Println("ID: " + mails[i].Id + " | Subject: " + mails[i].Subject)
		fmt.Println("")
	}
}

func DeleteFromArray(mails []mail.Mail, element mail.Mail) []mail.Mail {
	newMails := make([]mail.Mail, len(mails)-1)
	for i := 0; i < len(mails); i++ {
		if mails[i].Id != element.Id {
			newMails[i] = mails[i]
		}
	}
	return newMails
}

func GetAcknowledges(mails []mail.Mail) []mail.Acknowledgment {
	acknowledge := make([]mail.Acknowledgment, len(mails))
	for i := 0; i < len(mails); i++ {
		acknowledge[i].Id = mails[i].Id
		acknowledge[i].Subject = mails[i].Subject
	}
	return acknowledge
}
