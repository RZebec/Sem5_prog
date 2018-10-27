package main

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

	ticket, err := ticketContext.CreateNewTicket("TestTitle", ticket.Creator{Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"},
		ticket.MessageEntry{Id: 22, CreatorMail: "test@test.de", Content: "TestContent", OnlyInternal: false})
	fmt.Println(ticket)
	fmt.Println(err)
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

	if err := http.ListenAndServeTLS(":8080", "leaf.pem", "leaf.key", nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}
