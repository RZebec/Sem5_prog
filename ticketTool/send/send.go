package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
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


	// TODO: REMOVE only for Test
	/*apiClient, err := client.CreateClient(config)
	if err != nil {
		log.Fatal(err)
	}
	mails := *new([]mail.Mail)
	mails = append(mails, mail.Mail{Sender: "test@test.de", Subject: "testsubject", Content: "testcontent"} )
	err = apiClient.SendMails(mails)
	if err != nil {
		log.Fatal(err)
	}*/
	// TODO: Remove till here

	fmt.Print("write Email or load Email ? (w/l):")
	reader := bufio.NewReader(os.Stdin)
	message := make([]byte, 100)

	fmt.Print("write explicit mail or random mails ? (e/r):")
	email := mail.Mail{}
	eMails := []mail.Mail{}


	for true {
		answer := readEntry()
		if answer == "e" {
			eMails = email.ExplicitMail()
		} else if answer == "r" {
			validateNumberOfRandomMails(email, eMails)
		} else {
			fmt.Print("wrong entry. Please press e or r: ")
		}
	}

	fmt.Println("Email: ")
	fmt.Println("")
	fmt.Println("Start HTTPS-Request")

}
func readEntry() string {
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	value = strings.TrimRight(value, "\n")
	return value
}

func validateNumberOfRandomMails(email mail.Mail, eMails []mail.Mail) {
	for true {
		fmt.Println("Entry number of Random Mails: ")
		number, err := strconv.Atoi(readEntry())

		if err != nil {
			fmt.Println("Entry is no Number!")
		} else {
			eMails = email.RandomMail(number)
			break
		}
	}
}
