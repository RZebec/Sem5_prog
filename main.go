package main

import (
	"fmt"
	"html"
	"net/http"
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
	http.HandleFunc("/", foohandler)
	http.HandleFunc("/files/", tempHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}