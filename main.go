package main

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"./ticketsystem/webserver/core"
	"./ticketsystem/webserver/core/session"
	"./ticketsystem/webserver/webui"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func foohandler(w http.ResponseWriter, r *http.Request){
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
	//var loginHandler(logger, SessionManager)
	//staticFileHAndler := CreateNEw(config)
	//authenticationHandler
	config := config.Configuration{}
	filePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		config.LoginDataFolderPath = filePath
	} else {
		panic(err)
	}


	sessionManager := session.LoginSystem{}
	err = sessionManager.Initialize(path.Join(config.LoginDataFolderPath, "LoginData"))
	if err != nil {
		panic(err)
	}



	exampleHandler := webui.ExampleHtmlHandler{Prefix: "Das ist mein Prefix"}
	wrapper := core.Handler{Next: exampleHandler}



	http.HandleFunc("/", foohandler)
	http.HandleFunc("/files/", tempHandler)
	http.HandleFunc("/example", wrapper.ServeHTTP )

	if err := http.ListenAndServeTLS(":8080", "./ticketsystem/leaf.pem", "./ticketsystem/leaf.key", nil); err != nil {
		panic(err)
	}




	//staticFileHandlers.StaticFileHandler()
}