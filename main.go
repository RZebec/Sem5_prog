package main

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/login"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/logout"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/register"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func main() {
	// Core functionality
	// var logger = ...
	// var sessionmanager =

	//
	// interface logger ( LogDebug(), LogInfo())
	//
	// Website Handlers

	//staticFileHAndler := CreateNEw(config)
	//authenticationHandler
	config := config.Configuration{}
	filePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		config.LoginDataFolderPath = filePath
		config.TicketDataFolderPath = path.Join(filePath, "Tickets")
	} else {
		panic(err)
	}

	userContext := user.LoginSystem{}
	err = userContext.Initialize(path.Join(config.LoginDataFolderPath, "LoginData"))
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

	// TODO: Remove later, for test purposes only
	config.AccessTokenCookieName = "Access-Token"
	success, _ := userContext.Register("example@test.com", "1234", "Max", "Pimmelberg")
	if success {

	}

	indexPageHandler := webui.IndexPageHandler{UserContext: &userContext, Config: config}
	http.HandleFunc("/", indexPageHandler.ServeHTTP)

	filesHandler := files.FilesHandler{}
	http.HandleFunc("/files/", filesHandler.ServeHTTP)

	registerPageHandler := register.RegisterPageHandler{UserContext: &userContext, Config: config}
	http.HandleFunc("/register", registerPageHandler.ServeHTTP)

	registerHandler := register.RegisterHandler{UserContext: &userContext, Config: config}
	http.HandleFunc("/user_register", registerHandler.ServeHTTP)

	loginPageHandler := login.LoginPageHandler{UserContext: &userContext, Config: config}
	http.HandleFunc("/login", loginPageHandler.ServeHTTP)

	loginHandler := login.LoginHandler{UserContext: &userContext, Config: config}
	http.HandleFunc("/user_login", loginHandler.ServeHTTP)

	logoutHandler := logout.LogoutHandler{UserContext: &userContext, Config: config}
	logoutWrapper := wrappers.AuthenticationHandler{Next: logoutHandler, Config: config, UserContext: &userContext}
	http.HandleFunc("/user_logout", logoutWrapper.ServeHTTP)

	if err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}
