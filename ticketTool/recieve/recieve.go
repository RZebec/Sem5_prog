package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/saving"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/sharing"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
	"os"
)

func main() {
	logger := logging.ConsoleLogger{SetTimeStamp: true}
	config := configuration.Configuration{}
	config.RegisterFlags()
	config.BindFlags()

	if !config.ValidateConfiguration(logger) {
		fmt.Println("Configuration is not valid. Press enter to exit application.")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadByte()
		return
	}

	apiClient, createErr := client.CreateClient(config)
	if createErr != nil {
		panic(createErr)
	}

	mails := []mail.Mail{}
	for true {
		recieveMails, err := apiClient.ReceiveMails()
		if err != nil {
			fmt.Print("Transmission is going wrong. Retry? (n,press any key)")
			answer := inputOutput.ReadEntry()
			if answer == "n" {
				break
			}
		}
		mails = recieveMails
	}

	saving.SaveAcknowledge(sharing.GetAcknowledges(mails))

	sharing.SharingAllMails(mails)
	/*
	ackErr := apiClient.AcknowledgeMails(acknowledge)

	if ackErr != nil {
		fmt.Println("acknowlege is not posted")
	}
*/
}
