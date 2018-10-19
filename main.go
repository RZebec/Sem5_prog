package main

import (
	"net/http"
	"./staticFileHandlers"
)

func main() {
	// Core functionality
	// var logger = ...
	// var sessionmanager = 

	//
    // interface logger ( LogDebug(), LogInfo())
	//
	// Website Handlers
	//var loginHandler(logger, SessionManager)
	http.HandleFunc("/", staticFileHandlers.IndexHandler)
	http.HandleFunc("/login", staticFileHandlers.LoginPageHandler)
	http.HandleFunc("/files/styles", staticFileHandlers.CssHandler)
	http.HandleFunc("/files/login-styles", staticFileHandlers.LoginStyleHandler)
	http.HandleFunc("/files/javascript", staticFileHandlers.JsHandler)
	http.HandleFunc("/files/login", staticFileHandlers.LoginHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}