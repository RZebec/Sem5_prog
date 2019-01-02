package confirm

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
)

type Confirmation interface {
	GetAllAcknowledges(mails []mail.Mail) []mail.Acknowledgment
	GetSingleAcknowledges(allAcknowledges []mail.Acknowledgment, answer string) ([]mail.Acknowledgment, []mail.Acknowledgment)
	ShowAllEmailAcks(allAcknowledges []mail.Acknowledgment)
}

type Confirmator struct {

}

func (c *Confirmator) GetAllAcknowledges(mails []mail.Mail) []mail.Acknowledgment {
	acknowledge := make([]mail.Acknowledgment, len(mails))
	for i := 0; i < len(mails); i++ {
		acknowledge[i].Id = mails[i].Id
		acknowledge[i].Subject = mails[i].Subject
	}
	c.ShowAllEmailAcks(acknowledge)
	return acknowledge
}

func (c *Confirmator)  GetSingleAcknowledges(allAcknowledges []mail.Acknowledgment, answer string) ([]mail.Acknowledgment, []mail.Acknowledgment) {
	acknowledge := make([]mail.Acknowledgment, 1)
	for i := 0; i < len(allAcknowledges); i++ {
		if answer == allAcknowledges[i].Subject {
			acknowledge[0].Id = allAcknowledges[i].Id
			acknowledge[0].Subject = allAcknowledges[i].Subject
			newAcknowledges := c.deleteFromArray(allAcknowledges, allAcknowledges[i])
			return newAcknowledges, acknowledge
		}
	}
	fmt.Println("Subject not found")
	return allAcknowledges, make([]mail.Acknowledgment, 0)

}

func (c *Confirmator)  ShowAllEmailAcks(allAcknowledges []mail.Acknowledgment) {
	fmt.Println("All Emails: ")
	for i := 0; i < len(allAcknowledges); i++ {
		fmt.Println("ID: " + allAcknowledges[i].Id + " | Subject: " + allAcknowledges[i].Subject)
		fmt.Println("")
	}
}

func (c *Confirmator)  deleteFromArray(allAcknowledges []mail.Acknowledgment, element mail.Acknowledgment) []mail.Acknowledgment {
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
