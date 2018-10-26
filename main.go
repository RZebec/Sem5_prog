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

	exampleHandler := webui.ExampleHtmlHandler{Prefix: "Das ist mein Prefix"}
	wrapper := core.Handler{Next: exampleHandler}

	http.HandleFunc("/", foohandler)
	http.HandleFunc("/example", wrapper.ServeHTTP)

	filesHandler := webui.FilesHandler{}
	filesWrapper := core.Handler{Next: filesHandler}
	http.HandleFunc("/files/", filesWrapper.ServeHTTP)

	loginPageHandler := login.LoginPageHandler{IsUserLoggedIn: false, IsLoginFailed: false}
	loginPageWrapper := core.Handler{Next: loginPageHandler}
	http.HandleFunc("/login", loginPageWrapper.ServeHTTP)

	loginHandler := login.LoginHandler{UserManager: &sessionManager, LoginPageHandler: loginPageHandler}
	loginWrapper := core.Handler{Next: loginHandler}
	http.HandleFunc("/user_login", loginWrapper.ServeHTTP)

	logoutHandler := webui.LogoutHandler{}
	logoutWrapper := core.Handler{Next: logoutHandler}
	http.HandleFunc("/user_logout", logoutWrapper.ServeHTTP)

	if err := http.ListenAndServeTLS(":8080", "leaf.pem", "leaf.key", nil); err != nil {
		panic(err)
	}

	//staticFileHandlers.StaticFileHandler()
}
