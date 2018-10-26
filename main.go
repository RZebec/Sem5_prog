package main

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/login"
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

	sessionManager.Register("example@test.com", "1234")

	exampleHandler := webui.ExampleHtmlHandler{Prefix: "Das ist mein Prefix"}
	wrapper := core.Handler{Next: exampleHandler}

	indexPageHandler := webui.IndexPageHandler{}
	indexPageWrapper := core.Handler{Next: indexPageHandler}
	http.HandleFunc("/", indexPageWrapper.ServeHTTP)

	http.HandleFunc("/example", wrapper.ServeHTTP)

	filesHandler := webui.FilesHandler{}
	http.HandleFunc("/files/", filesHandler.ServeHTTP)

	loginPageHandler := login.LoginPageHandler{IsUserLoggedIn: false, IsLoginFailed: false}
	http.HandleFunc("/login", loginPageHandler.ServeHTTP)

	loginHandler := login.LoginHandler{UserManager: &sessionManager}
	http.HandleFunc("/user_login", loginHandler.ServeHTTP)

	logoutHandler := webui.LogoutHandler{UserManager: &sessionManager}
	logoutWrapper := core.Handler{Next: logoutHandler}
	http.HandleFunc("/user_logout", logoutWrapper.ServeHTTP)

	if err := http.ListenAndServeTLS(":8080", "leaf.pem", "leaf.key", nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}
