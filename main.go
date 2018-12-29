package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/api"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/api/mails"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
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

	if !configuration.ValidateConfiguration(logger) {
		fmt.Println("Configuration is not valid. Press enter to exit application.")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadByte()
		return
	}

	mailContext := mail.MailManager{}
	err = mailContext.Initialize(configuration.MailDataFolderPath, configuration.SendingMailAddress, logger)
	mailContext.CreateNewOutgoingMail("test1@test1.de", "testSubject1", "TestContent1")
	mailContext.CreateNewOutgoingMail("test2@test2.de", "testSubject2", "TestContent2")
	mailContext.CreateNewOutgoingMail("test3@test2.de", "testSubject3", "TestContent2")

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

	http.HandleFunc("/api/mail/incoming", getIncomingMailHandlerChain(*apiConfig, &mailContext, logger).ServeHTTP)

	handlerManager := webui.HandlerManager{
		UserContext:   &userContext,
		TicketContext: &ticketContext,
		Config:        configuration,
		Logger:        logger,
	}

	templateManager.LoadTemplates(logger)
	handlerManager.StartServices()

	if err := http.ListenAndServeTLS(configuration.GetServiceUrl(), configuration.CertificatePath, configuration.CertificateKeyPath, nil); err != nil {
		logger.LogError("Main", err)
	}

	//staticFileHandlers.StaticFileHandler()
}

func getIncomingMailHandlerChain(apiConfig config.ApiConfiguration, mailContext mail.MailContext, logger logging.Logger) http.Handler {
	incomingMailHandler := mails.IncomingMailHandler{Logger: logger, MailContext: mailContext}
	apiAuthenticationHandler := api.ApiKeyAuthenticationHandler{ApiKeyResolver: apiConfig.GetIncomingMailApiKey,
		Next: &incomingMailHandler, AllowedMethod: "POST", Logger: logger}
	return &apiAuthenticationHandler
}

func getAcknowledgeMailHandlerChain(apiConfig config.ApiConfiguration, mailContext mail.MailContext, logger logging.Logger) http.Handler {
	incomingMailHandler := mails.AcknowledgeMailHandler{Logger: logger, MailContext: mailContext}
	apiAuthenticationHandler := api.ApiKeyAuthenticationHandler{ApiKeyResolver: apiConfig.GetIncomingMailApiKey,
		Next: &incomingMailHandler, AllowedMethod: "POST", Logger: logger}
	return &apiAuthenticationHandler
}

func getOutgoingMailHandlerChain(apiConfig config.ApiConfiguration, mailContext mail.MailContext, logger logging.Logger) http.Handler {
	outgoingMailHandler := mails.OutgoingMailHandler{Logger: logger, MailContext: mailContext}
	apiAuthenticationHandler := api.ApiKeyAuthenticationHandler{ApiKeyResolver: apiConfig.GetOutgoingMailApiKey,
		Next: &outgoingMailHandler, AllowedMethod: "GET", Logger: logger}
	return &apiAuthenticationHandler
}
