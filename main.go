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

	testMessageContent := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed ultrices congue mi eget consequat. Nulla facilisi. 
						Maecenas bibendum diam eget lectus sagittis mollis.Fusce vel tempus nisl. Quisque sollicitudin ultrices tristique. 
						Integer a nisi vitae justo hendrerit facilisis. Integer a eros eu erat euismod tristique vitae vitae orci. 
						Nunc tincidunt faucibus turpis, sed cursus tellus dapibus ut. Aliquam metus ligula, elementum eu auctor at, rutrum at nulla. 
						Mauris sit amet mauris vel velit congue rhoncus. Donec a dolor luctus, mattis ligula vitae, elementum libero. Morbi pellentesque scelerisque suscipit. 
						Etiam sit amet tincidunt ex. Vivamus quis magna ornare, elementum arcu ac, rhoncus tellus. Nullam ullamcorper pharetra sodales. 
						Nunc purus nibh, vestibulum quis ex non, ornare congue nibh. Nulla sagittis magna aliquet malesuada gravida. `

	testMessages := []ticket.MessageEntry{
		{Id: 0, CreatorMail: "test@test.de", Content: testMessageContent, OnlyInternal: false},
		{Id: 1, CreatorMail: "test1@test.de", Content: testMessageContent, OnlyInternal: false},
		{Id: 2, CreatorMail: "test2@test.de", Content: testMessageContent, OnlyInternal: true},
		{Id: 3, CreatorMail: "test3@test.de", Content: testMessageContent, OnlyInternal: false},
	}


	ticketg, err := ticketContext.CreateNewTicket("TestTitle", ticket.Creator{Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"}, testMessages[0])
	ticketContext.AppendMessageToTicket(1, testMessages[1])
	ticketContext.AppendMessageToTicket(1, testMessages[3])
	fmt.Println(ticketg)

	ticketg, err = ticketContext.CreateNewTicket("TestTitle2", ticket.Creator{Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"}, testMessages[0])
	ticketContext.AppendMessageToTicket(2, testMessages[1])
	ticketContext.AppendMessageToTicket(2, testMessages[3])
	fmt.Println(ticketg)

	ticketg, err = ticketContext.CreateNewTicketForInternalUser("TestTitle", user.User{UserId: 1, Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"}, testMessages[2])
	fmt.Println(ticketg)

	ticketg, err = ticketContext.CreateNewTicketForInternalUser("TestTitle", user.User{UserId: 2, Mail: "peter@test.de", FirstName: "Peter", LastName: "Test"}, testMessages[2])
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
		TicketContext: &ticketContext,
		Config: config,
	}

	handlerManager.StartServices()

	if err := http.ListenAndServeTLS(config.GetServiceUrl(), config.CertificatePath, config.CertificateKeyPath, nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}
