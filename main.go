package main

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/login"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/logout"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
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

	// TODO: Remove later, for test purposes only
	config.AccessTokenCookie = helpers.Cookie{Name: "Access-Token"}
	sessionManager.Register("example@test.com", "1234")

	indexPageHandler := webui.IndexPageHandler{SessionManager: &sessionManager, AccessTokenCookie: config.AccessTokenCookie}
	http.HandleFunc("/", indexPageHandler.ServeHTTP)

	filesHandler := files.FilesHandler{}
	http.HandleFunc("/files/", filesHandler.ServeHTTP)

	loginPageHandler := login.LoginPageHandler{SessionManager: &sessionManager, AccessTokenCookie: config.AccessTokenCookie}
	http.HandleFunc("/login", loginPageHandler.ServeHTTP)

	loginHandler := login.LoginHandler{UserManager: &sessionManager, AccessTokenCookie: config.AccessTokenCookie}
	http.HandleFunc("/user_login", loginHandler.ServeHTTP)

	logoutHandler := logout.LogoutHandler{UserManager: &sessionManager, AccessTokenCookie: config.AccessTokenCookie}
	logoutWrapper := wrappers.AuthenticationHandler{Next: logoutHandler, AccessTokenCookie: config.AccessTokenCookie, SessionManager: &sessionManager}
	http.HandleFunc("/user_logout", logoutWrapper.ServeHTTP)

	if err := http.ListenAndServeTLS(":8080", "leaf.pem", "leaf.key", nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}
