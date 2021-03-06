// 5894619, 6720876, 9793350
package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/acknowledgementStorage"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/confirm"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/recviever"
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

	io := inputOutput.DefaultInputOutput{}
	confirmator := confirm.Confirmator{}

	recieve := recviever.CreateReciever(config, &io, &apiClient, storage, &confirmator)

	fmt.Println("Recieve Mails")
	for true {
		ack, result := recieve.Run()
		if ack == false {
			break
		}
		if ack == true && result != nil {
			fmt.Println(result)
		}
		if ack == false && result != nil {
			break
		}
	}

}
