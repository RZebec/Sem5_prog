package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"fmt"
	"net/http"
	"os"
)

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

	templateManager.LoadTemplates()

	// TODO: Remove later, for test purposes only
	config.AccessTokenCookieName = "Access-Token"
	success, _ := userContext.Register("example@test.com", "1234", "Max", "Pimmelberg")
	if success {

	}

	handlerManager := webui.HandlerManager{
		UserContext: &userContext,
		Config: config,
	}

	handlerManager.StartServices()

	if err := http.ListenAndServeTLS(config.GetServiceUrl(), config.CertificatePath, config.CertificateKeyPath, nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}
