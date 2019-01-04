package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/acknowledgementStorage"
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

	// Create api client:
	apiClient, createErr := client.CreateClient(config)
	if createErr != nil {
		panic(createErr)
	}

	// Create acknowledge storage:
	storage, createErr := acknowledgementStorage.InitializeAckManager(config.UnAcknowledgedMailPath)
	if createErr != nil {
		panic(createErr)
	}
	// TODO Remove the following line: Only to enable compilation till it is used
	fmt.Println(storage)

	fmt.Println("Recieve Mails")
	mails := []mail.Mail{}
	for true {
		recieveMails, err := apiClient.ReceiveMails()
		if err != nil {
			fmt.Print("Transmission is going wrong. Retry? (n,press any key)")
			answer := inputOutput.ReadEntry()
			if answer == "n" {
				break
			}
		} else {
			fmt.Println("Mails are incoming")
			mails = recieveMails
			allOrSpecifySharing(apiClient, &mails)
			break
		}
	}

}

func allOrSpecifySharing(apiClient client.ApiClient, mails *[]mail.Mail) {
	//acknowledgement.WriteAcknowledgements(sharing.GetAcknowledges(mails))

	for true {
		fmt.Println("share all Messages or specify Messages? (all/specify):")
		answer := inputOutput.ReadEntry()
		if answer == "all" {
			acknowledge := sharing.ShareAllMails(*mails)
			ackError := apiClient.AcknowledgeMails(acknowledge)
			if ackError != nil {
				fmt.Println("acknowlege is not posted")
			} else {
				//acknowledgement.DeleteAcknowledges(acknowledge)
				break
			}
		} else if answer == "specify" {
			acknowledge, newMails := sharing.ShareSingleMails(*mails)
			mails = &newMails
			ackError := apiClient.AcknowledgeMails(acknowledge)
			if ackError != nil {
				fmt.Println("acknowlege is not posted")
			} else {
				//acknowledgement.DeleteAcknowledges(acknowledge)
			}
			if len(*mails) == 0 {
				break
			}
		}
	}
}
