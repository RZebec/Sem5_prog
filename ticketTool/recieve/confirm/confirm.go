package confirm

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
)

func GetAllAcknowledges(mails []mail.Mail) []mail.Acknowledgment {
	acknowledge := make([]mail.Acknowledgment, len(mails))
	for i := 0; i < len(mails); i++ {
		acknowledge[i].Id = mails[i].Id
		acknowledge[i].Subject = mails[i].Subject
	}
	showAllEmailAcks(acknowledge)
	return acknowledge
}

func GetSingleAcknowledges(allAcknowledges []mail.Acknowledgment) ([]mail.Acknowledgment, []mail.Acknowledgment) {
	acknowledge := make([]mail.Acknowledgment, 1)
	showAllEmailAcks(allAcknowledges)
	fmt.Println("Specify Acknowledge by Subject: ")
	answer := inputOutput.ReadEntry()
	for i := 0; i < len(allAcknowledges); i++ {
		if answer == allAcknowledges[i].Subject {
			acknowledge[0].Id = allAcknowledges[i].Id
			acknowledge[0].Subject = allAcknowledges[i].Subject
			newAcknowledges := deleteFromArray(allAcknowledges, allAcknowledges[i])
			return newAcknowledges, acknowledge
		}
	}
	fmt.Println("Subject not found")
	return allAcknowledges, acknowledge

}

func showAllEmailAcks(allAcknowledges []mail.Acknowledgment) {
	fmt.Println("All Emails: ")
	for i := 0; i < len(allAcknowledges); i++ {
		fmt.Println("ID: " + allAcknowledges[i].Id + " | Subject: " + allAcknowledges[i].Subject)
		fmt.Println("")
	}
}

func deleteFromArray(allAcknowledges []mail.Acknowledgment, element mail.Acknowledgment) []mail.Acknowledgment {
	newAcknowledges := make([]mail.Acknowledgment, len(allAcknowledges)-1)
	j := 0
	for i := 0; i < len(allAcknowledges); i++ {
		if allAcknowledges[i].Id != element.Id {
			newAcknowledges[j] = allAcknowledges[i]
			j++
		}
	}
	return newAcknowledges
}
