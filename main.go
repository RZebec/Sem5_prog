package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/api"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/api/mails"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
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


func main() {
	logger := logging.ConsoleLogger{SetTimeStamp: true}
	configuration := config.Configuration{}
	configuration.RegisterFlags()
	configuration.BindFlags()

	apiConfig, err := config.CreateAndInitialize(configuration)
	if err != nil {
		panic(err)
	}
	fmt.Println(apiConfig)

	if !configuration.ValidateConfiguration(logger) {
		fmt.Println("Configuration is not valid. Press enter to exit application.")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadByte()
		return
	}

	userContext := user.LoginSystem{}
	err = userContext.Initialize(configuration.LoginDataFolderPath)
	if err != nil {
		panic(err)
	}

	ticketContext := ticket.TicketManager{}
	ticketContext.Initialize(configuration.TicketDataFolderPath)

	ticketa, err := ticketContext.CreateNewTicket("TestTitle", ticket.Creator{Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"},
		ticket.MessageEntry{Id: 0, CreatorMail: "test@test.de", Content: "TestContent1", OnlyInternal: false})
	fmt.Println(ticketa)
	ticketa, err = ticketContext.CreateNewTicket("TestTitle2", ticket.Creator{Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"},
		ticket.MessageEntry{Id: 0, CreatorMail: "test@test.de", Content: "TestContent2", OnlyInternal: false})

	ticketg, err := ticketContext.CreateNewTicketForInternalUser("TestTitle", user.User{UserId: 1, Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"},
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
	http.HandleFunc("/api/mail/incoming", getIncomingMailHandlerChain(*apiConfig).ServeHTTP)

	if err := http.ListenAndServeTLS(configuration.GetServiceUrl(), configuration.CertificatePath, configuration.CertificateKeyPath, nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}

func getIncomingMailHandlerChain(apiConfig config.ApiConfiguration) http.Handler {
	incomingMailHandler := mails.IncomingMailHandler{}
	apiAuthenticationHandler := api.ApiKeyAuthenticationHandler{ApiKeyResolver: apiConfig.GetIncomingMailApiKey,
		Next: &incomingMailHandler, AllowedMethod: "POST"}
	return &apiAuthenticationHandler
}
