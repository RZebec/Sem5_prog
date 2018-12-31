package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/acknowledgementStorage"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"fmt"
	"os"
	"testing"
)

func TestReciever(t *testing.T) {
	logger := logging.ConsoleLogger{SetTimeStamp: true}
	config := configuration.Configuration{}
	config.RegisterFlags()
	config.BindFlags()
	config.CertificatePath = "test_cert.pem"
	config.ApiKeysFilePath = "test_api.keys"

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

	/*
		without Server
	*/
	err := reciever(apiClient, storage)
	fmt.Println(err)
}
