package sender

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/mailGeneration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"strconv"
)

type Sender struct {
	io inputOutput.InputOutput
	sendConfig configuration.Configuration
	apiClient client.Client
	mailGenerator mailGeneration.MailGeneration
}

func CreateSender(config configuration.Configuration, io inputOutput.InputOutput,
		apiClient client.Client, mailGen mailGeneration.MailGeneration) Sender{

	sender := Sender{io: io, sendConfig: config, apiClient: apiClient, mailGenerator: mailGen}
	return sender
}

func (s *Sender) Run(){
	s.io.Print("write explicit mail or random mails ? (e/r):")
	var eMails []mail.Mail
	for true {
		answer := s.io.ReadEntry()
		if answer == "e" {
			eMails = s.mailGenerator.ExplicitMail()
			s.httpRequest(eMails)
			break
		} else if answer == "r" {
			foundNumber := 0
			for true {
				number, err := s.entryNumberOfRandomMails()
				if err == nil {
					foundNumber = number
					break
				}
			}
			eMails = s.mailGenerator.RandomMail(foundNumber, 10, 50)
			s.httpRequest(eMails)
			break
		} else {
			s.io.Print("wrong entry. Please press e or r: ")
		}
	}
}

func (s *Sender) httpRequest(eMails []mail.Mail) {
	s.io.Print("Start HTTPS-Request")
	err := s.apiClient.SendMails(eMails)
	if err != nil {
		s.io.Print(err.Error())
	}
}

func  (s *Sender) entryNumberOfRandomMails() (int, error) {
	s.io.Print("Entry number of Random Mails: ")
	number, err := strconv.Atoi(s.io.ReadEntry())
	if err != nil {
		s.io.Print("Entry is no Number!")
		return 0, err
	} else {
		return number, nil
	}
}