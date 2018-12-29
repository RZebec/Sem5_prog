package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"

	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
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

	apiClient, err := client.CreateClient(config)
	if err != nil {
		panic(err)
	}
	receivedMails, err := apiClient.ReceiveMails()
	if err == nil {
		fmt.Println(receivedMails)
		var acks []mail.Acknowledgment
		for _, mailToAck := range receivedMails{
			acks = append(acks, mail.Acknowledgment{Id: mailToAck.Id})
		}
		err = apiClient.AcknowledgeMails(acks)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

	//message := clientContainer.HttpsRequest(config.BaseUrl, config.Port, config.CertificatePath, "Test")
	//fmt.Println(message)
}
