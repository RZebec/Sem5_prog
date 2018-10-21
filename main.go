package main

import (
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
	staticFileHandlers.StaticFileHandler()
}