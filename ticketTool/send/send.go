package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	fmt.Print("write explicit mail or random mails ? (e/r):")

	email := mail.Mail{}
	eMails := []mail.Mail{}

	for true {
		answer := readEntry()
		if answer == "e" {
			eMails = email.ExplicitMail()
			httpRequest(config, eMails)
			break
		} else if answer == "r" {
			number := entryNumberOfRandomMails()
			eMails = email.RandomMail(number)
			httpRequest(config, eMails)
			break
		} else {
			fmt.Print("wrong entry. Please press e or r: ")
		}
	}

}

func httpRequest(config configuration.Configuration, eMails []mail.Mail) {
	fmt.Println("Start HTTPS-Request")
	apiClient, err := client.CreateClient(config)
	if err != nil {
		log.Fatal(err)
	}
	apiClient.SendMails(eMails)
}

func readEntry() string {
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	value = strings.TrimRight(value, "\n")
	return value
}

func entryNumberOfRandomMails() int {
	for true {
		fmt.Println("Entry number of Random Mails: ")
		number, err := strconv.Atoi(readEntry())
		if err != nil {
			fmt.Println("Entry is no Number!")
		} else {
			return number
		}
	}
}
