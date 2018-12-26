package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputContainer"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"fmt"
	"os"
)

func main()  {
	logger := logging.ConsoleLogger{SetTimeStamp: true}
	container := inputContainer.Configuration{}
	container.RegisterFlags()
	container.BindFlags()

	if !container.ValidateConfiguration(logger) {
		fmt.Println("Configuration is not valid. Press enter to exit application.")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadByte()
		return
	}

	
}
