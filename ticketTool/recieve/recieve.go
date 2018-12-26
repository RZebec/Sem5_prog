package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"fmt"
	"os"
)

func main()  {
	logger := logging.ConsoleLogger{SetTimeStamp: true}
	config := config.Configuration{}
	config.RegisterFlags()
	config.BindFlags()

	if !config.ValidateConfiguration(logger) {
		fmt.Println("Configuration is not valid. Press enter to exit application.")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadByte()
		return
	}
}
