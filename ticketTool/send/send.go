package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/mailGeneration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	logger := logging.ConsoleLogger{SetTimeStamp: true}
	config := configuration.Configuration{}
	config.RegisterFlags()
	config.BindFlags()

	mailGenerator := mailGeneration.MailGenerator{}

	if !config.ValidateConfiguration(logger) {
		fmt.Println("Configuration is not valid. Press enter to exit application.")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadByte()
		return
	}
	fmt.Print("write explicit mail or random mails ? (e/r):")

	var eMails []mail.Mail

	for true {
		answer := inputOutput.ReadEntry()
		if answer == "e" {
			eMails = mailGenerator.ExplicitMail()
			httpRequest(config, eMails)
			break
		} else if answer == "r" {
			number := entryNumberOfRandomMails()
			eMails = mailGenerator.RandomMail(number)
			httpRequest(config, eMails)
			break
		} else {
			fmt.Print("wrong entry. Please press e or r: ")
		}
	}
}

func httpRequest(config configuration.Configuration, eMails []mail.Mail) {
	fmt.Println("Start HTTPS-Request")
	apiClient, err := client.CreateClient(config)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Add error handling
	apiClient.SendMails(eMails)
	// TODO REMOVE ONLY FOR TEST
	var acks []mail.Acknowledgment
	for _, email := range eMails {
		acks = append(acks, mail.Acknowledgment{Id: email.Id, Subject: email.Subject})
	}
	apiClient.AcknowledgeMails(acks)
	// TODO REMOVE TILL HERE
}

func entryNumberOfRandomMails() int {
	for true {
		fmt.Println("Entry number of Random Mails: ")
		number, err := strconv.Atoi(inputOutput.ReadEntry())
		if err != nil {
			fmt.Println("Entry is no Number!")
		} else {
			return number
		}
	}
	return 0
}
