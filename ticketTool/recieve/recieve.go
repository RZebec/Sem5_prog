package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"

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
	fmt.Println(apiClient.ReceiveMails())
	//message := clientContainer.HttpsRequest(config.BaseUrl, config.Port, config.CertificatePath, "Test")
	//fmt.Println(message)
}
