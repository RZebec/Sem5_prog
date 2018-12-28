package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func foohandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	w.Write([]byte("HHH"))
}

func tempHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	w.Write([]byte(r.URL.Path))
}

func incomingApiHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t []mail.Mail
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t)
	}
	w.WriteHeader(200)
}

func main() {
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

	userContext := user.LoginSystem{}
	err := userContext.Initialize(config.LoginDataFolderPath)
	if err != nil {
		panic(err)
	}

	ticketContext := ticket.TicketManager{}
	ticketContext.Initialize(config.TicketDataFolderPath)

	ticketg, err := ticketContext.CreateNewTicket("TestTitle", ticket.Creator{Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"},
		ticket.MessageEntry{Id: 0, CreatorMail: "test@test.de", Content: "TestContent", OnlyInternal: false})
	fmt.Println(ticketg)
	ticketg, err = ticketContext.CreateNewTicket("TestTitle2", ticket.Creator{Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"},
		ticket.MessageEntry{Id: 0, CreatorMail: "test@test.de", Content: "TestContent", OnlyInternal: false})

	ticketg, err = ticketContext.CreateNewTicketForInternalUser("TestTitle", user.User{UserId: 1, Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"},
		ticket.MessageEntry{Id: 0, CreatorMail: "test@test.de", Content: "TestContent", OnlyInternal: false})
	fmt.Println(ticketg)

	ticketg, err = ticketContext.CreateNewTicketForInternalUser("TestTitle", user.User{UserId: 2, Mail: "peter@test.de", FirstName: "Peter", LastName: "Test"},
		ticket.MessageEntry{Id: 0, CreatorMail: "test@test.de", Content: "TestContent", OnlyInternal: true})
	fmt.Println(ticketg)

	exists, ticket := ticketContext.GetTicketById(2)
	fmt.Println(exists)
	fmt.Println(ticket)

	g := ticketContext.GetAllTicketInfo()
	fmt.Println(len(g))

	exampleHandler := webui.ExampleHtmlHandler{Prefix: "Das ist mein Prefix"}
	wrapper := core.Handler{Next: exampleHandler}

	http.HandleFunc("/", foohandler)
	http.HandleFunc("/files/", tempHandler)
	http.HandleFunc("/example", wrapper.ServeHTTP)
	http.HandleFunc("/api/mail/incoming", incomingApiHandler)

	if err := http.ListenAndServeTLS(config.GetServiceUrl(), config.CertificatePath, config.CertificateKeyPath, nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}
