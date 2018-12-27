package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketTool/clientContainer"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputContainer"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

	fmt.Print("write Email or load Email ? (w/l):")
	reader := bufio.NewReader(os.Stdin)
	message := make([]byte, 100)

	for true {
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimRight(answer, "\n")
		if answer == "w" {
			writeEmail()
			break
		} else if answer == "l" {
			for true {
				fmt.Print("Enter email path:")
				path := readEntry()
				answer, available := loadEmail(path)
				if available == true {
					message = answer
					break
				}
			}
		} else {
			fmt.Print("wrong entry. Please press w or l: ")
		}
	}

	fmt.Println("Email: ")
	fmt.Println(string(message))
	fmt.Println("")
	fmt.Println("Start HTTPS-Request")

	clientContainer.HttpsRequest(inputContainer.BaseUrl, inputContainer.Port, inputContainer.CertificatePath, string(message))
}

func writeEmail() {

}

func readEntry() string {
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	value = strings.TrimRight(value, "\n")
	return value
}

func loadEmail(path string) ([]byte, bool) {
	logger := logging.ConsoleLogger{SetTimeStamp: true}
	exists, err := helpers.FilePathExists(path)
	if err != nil {
		logger.LogError("Load File failure", err)
		return make([]byte, 0), false
	}
	if exists == false {
		fmt.Println("File doesnt exist")
		return make([]byte, 0), false
	} else {
		return loadFile(path)
	}
}

func loadFile(path string) ([]byte, bool) {
	logger := logging.ConsoleLogger{SetTimeStamp: true}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		logger.LogError("Read File is going wrong", err)
		return dat, false
	}
	return dat, true
}
