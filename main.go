package main

import (
	"bufio"
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/api"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/api/mails"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

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

	mailContext := mailData.MailManager{}
	err = mailContext.Initialize(configuration.MailDataFolderPath, configuration.SendingMailAddress, logger)

	userContext := userData.LoginSystem{}
	err = userContext.Initialize(configuration.LoginDataFolderPath)
	if err != nil {
		panic(err)
	}

	ticketContext := ticketData.TicketManager{}
	ticketContext.Initialize(configuration.TicketDataFolderPath)

	http.HandleFunc(shared.SendPath, getIncomingMailHandlerChain(*apiConfig, &mailContext, &ticketContext, &userContext, logger).ServeHTTP)
	http.HandleFunc(shared.AcknowledgmentPath, getAcknowledgeMailHandlerChain(*apiConfig, &mailContext, logger).ServeHTTP)
	http.HandleFunc(shared.ReceivePath, getOutgoingMailHandlerChain(*apiConfig, &mailContext, logger).ServeHTTP)

	templateMan := templateManager.TemplateManager{Templates: map[string]*template.Template{}}

	err = templateMan.LoadTemplates(logger)
	if err != nil {
		panic(err)
	}

	handlerManager := webui.HandlerManager{
		UserContext:      &userContext,
		TicketContext:    &ticketContext,
		Config:           configuration,
		Logger:           logger,
		ApiConfiguration: apiConfig,
		TemplateManager:  &templateMan,
		MailContext:      &mailContext,
	}

	templateMan.LoadTemplates(logger)
	handlerManager.RegisterHandlers()

	logger.LogInfo("Server", "Server started")
	if err := http.ListenAndServeTLS(configuration.GetServiceUrl(), configuration.CertificatePath, configuration.CertificateKeyPath, nil); err != nil {
		logger.LogError("Main", err)
	}
}

/*
	Get the api handler chain for incoming mails:
 */
func getIncomingMailHandlerChain(apiConfig config.ApiConfiguration, mailContext mailData.MailContext, ticketContext ticketData.TicketContext,
	userContext userData.UserContext, logger logging.Logger) http.Handler {
	incomingMailHandler := mails.IncomingMailHandler{Logger: logger, MailContext: mailContext, TicketContext: ticketContext,
		UserContext: userContext, MailRepliesFilter: &mails.RepliesFilter{}}
	apiAuthenticationHandler := api.ApiKeyAuthenticationHandler{ApiKeyResolver: apiConfig.GetIncomingMailApiKey,
		Next: &incomingMailHandler, AllowedMethod: "POST", Logger: logger}
	return &apiAuthenticationHandler
}

/*
	Get the api handler chain for acknowledgment of mails:
 */
func getAcknowledgeMailHandlerChain(apiConfig config.ApiConfiguration, mailContext mailData.MailContext, logger logging.Logger) http.Handler {
	incomingMailHandler := mails.AcknowledgeMailHandler{Logger: logger, MailContext: mailContext}
	apiAuthenticationHandler := api.ApiKeyAuthenticationHandler{ApiKeyResolver: apiConfig.GetIncomingMailApiKey,
		Next: &incomingMailHandler, AllowedMethod: "POST", Logger: logger}
	return &apiAuthenticationHandler
}

/*
	Get the api handler chain for outgoing mails:
 */
func getOutgoingMailHandlerChain(apiConfig config.ApiConfiguration, mailContext mailData.MailContext, logger logging.Logger) http.Handler {
	outgoingMailHandler := mails.OutgoingMailHandler{Logger: logger, MailContext: mailContext}
	apiAuthenticationHandler := api.ApiKeyAuthenticationHandler{ApiKeyResolver: apiConfig.GetOutgoingMailApiKey,
		Next: &outgoingMailHandler, AllowedMethod: "GET", Logger: logger}
	return &apiAuthenticationHandler
}
