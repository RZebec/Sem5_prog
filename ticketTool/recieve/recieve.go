package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/acknowledgement"
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

	acknowledgement.SaveAcknowledges(sharing.GetAcknowledges(mails))

	for true {
		fmt.Println("share all Messages or specify Messages? (all/specify)")
		answer := inputOutput.ReadEntry()
		if answer == "all" {
			acknowledge := sharing.ShareAllMails(mails)
			ackError := apiClient.AcknowledgeMails(acknowledge)
			if ackError != nil {
				fmt.Println("acknowlege is not posted")
			} else {
				acknowledgement.DeleteAcknowledges(acknowledge)
				break
			}
		} else if answer == "specify" {
			acknowledge := sharing.ShareSingleMails(mails)
			ackError := apiClient.AcknowledgeMails(acknowledge)
			if ackError != nil {
				fmt.Println("acknowlege is not posted")
			} else {
				acknowledgement.DeleteAcknowledges(acknowledge)
				break
			}
		}
	}

}
