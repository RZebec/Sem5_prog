package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/acknowledgementStorage"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/confirm"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
	"os"
	"strconv"
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

	fmt.Println("Recieve Mails")
	for true {
		reciever(apiClient, storage)
	}
}

func allOrSpecifyConfirm(apiClient client.ApiClient, allAcknowledges *[]mail.Acknowledgment, storage acknowledgementStorage.AckStorage) {

	for true {
		fmt.Println("send all Acknowledges or specify Acknowledges to Server. Or stop reciever (all/specify/stop):")
		answer := inputOutput.ReadEntry()
		if answer == "all" {
			ackError := apiClient.AcknowledgeMails(*allAcknowledges)
			if ackError != nil {
				fmt.Println("acknowlege is not posted")
			} else {
				fmt.Println("E-Mails are Acknowledged: ")
				storage.DeleteAcknowledges(*allAcknowledges)
				break
			}
		} else if answer == "specify" {
			confirm.ShowAllEmailAcks(*allAcknowledges)
			fmt.Println("Specify Acknowledge by Subject: ")
			answer := inputOutput.ReadEntry()
			newAcknowledges, selectedAck := confirm.GetSingleAcknowledges(*allAcknowledges, answer)
			allAcknowledges = &newAcknowledges
			ackError := apiClient.AcknowledgeMails(selectedAck)
			if ackError != nil {
				fmt.Println("acknowlege is not posted")
			} else {
				storage.DeleteAcknowledges(selectedAck)
				fmt.Println("E-Mail is Acknowledged: ")
			}
			if len(*allAcknowledges) == 0 {
				break
			}
		} else if answer == "stop" {
			break
		}
	}
}

func reciever(apiClient client.ApiClient, storage acknowledgementStorage.AckStorage) error {
	recieveMails, err := apiClient.ReceiveMails()
	if err != nil {
		fmt.Print("Transmission is going wrong. Retry? (n,press any key)")
		answer := inputOutput.ReadEntry()
		if answer == "n" {
			return err
		}
	} else {
		fmt.Println(strconv.Itoa(len(recieveMails)) + " Mails are coming from Server")
		acknowledges := confirm.GetAllAcknowledges(recieveMails)
		storage.AppendAcknowledgements(acknowledges)
		fmt.Println("Save Acknowledges...")
		allAcknowledges, err := storage.ReadAcknowledgements()
		if err != nil {
			fmt.Println("couldn't read storaged Acknowledges")
			return err
		} else if len(allAcknowledges) == 0 {
			fmt.Println("No Emails available")
			return nil
		}
		fmt.Println("Available Mails: " + strconv.Itoa(len(allAcknowledges)))
		if len(allAcknowledges) != 0 {
			allOrSpecifyConfirm(apiClient, &allAcknowledges, storage)
			return nil
		}
	}
	return nil
}
