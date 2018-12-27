package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/clientContainer"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputContainer"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"fmt"
	"os"
)

func main() {
	logger := logging.ConsoleLogger{SetTimeStamp: true}
	inputContainer := inputContainer.Configuration{}
	inputContainer.RegisterFlags()
	inputContainer.BindFlags()

	if !inputContainer.ValidateConfiguration(logger) {
		fmt.Println("Configuration is not valid. Press enter to exit application.")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadByte()
		return
	}

	clientContainer.HttpRequest(inputContainer.BaseUrl, inputContainer.Port, inputContainer.CertificatePath, "Test")
}
