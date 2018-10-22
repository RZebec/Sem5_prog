package main

import (
	"fmt"
	"html"
	"net/http"
	"./ticketsystem/webserver/webui"
	"./ticketsystem/webserver/core"

)

func foohandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Hello")
	fmt.Println("Hello, %q", html.EscapeString(r.URL.Path))
	w.Write([]byte("HHH"))
	}


func tempHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	fmt.Println("Hello Files, %q", html.EscapeString(r.URL.Path))
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

	exampleHandler := webui.ExampleHtmlHandler{Prefix: "Das ist mein Prefix"}
	wrapper := core.Handler{Next: exampleHandler}

	http.HandleFunc("/", foohandler)
	http.HandleFunc("/files/", tempHandler)
	http.HandleFunc("/example", wrapper.ServeHTTP )

	if err := http.ListenAndServeTLS(":8080", "./ticketsystem/92860317_localhost.cert", "./ticketsystem/92860317_localhost.key", nil); err != nil {
		panic(err)
	}




	//staticFileHandlers.StaticFileHandler()
}