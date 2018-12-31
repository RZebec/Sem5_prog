package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/mailGeneration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/send/sender"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"fmt"
	"log"
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

	io := inputOutput.DefaultInputOutput{}
	mailGenerator := mailGeneration.CreateMailGenerator(&io)

	apiClient, err := client.CreateClient(config)
	if err != nil  {
		log.Fatal(err)
	}

	send := sender.CreateSender(config, &io, &apiClient, &mailGenerator)
	send.Run()




}


