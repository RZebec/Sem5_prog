package webui

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/index"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/login"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/logout"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/register"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/ticketView"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/ticketexplorer"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
)

type HandlerManager struct {
	UserContext user.UserContext
	TicketContext ticket.TicketContext
	Config config.Configuration
}

func (handlerManager *HandlerManager) StartServices() {
	indexPageHandler := index.IndexPageHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config}
	http.HandleFunc("/", indexPageHandler.ServeHTTP)

	filesHandler := files.FilesHandler{}
	http.HandleFunc("/files/", filesHandler.ServeHTTP)

	registerPageHandler := register.RegisterPageHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config}
	http.HandleFunc("/register", registerPageHandler.ServeHTTP)

	registerHandler := register.RegisterHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config}
	http.HandleFunc("/user_register", registerHandler.ServeHTTP)

	loginPageHandler := login.LoginPageHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config}
	http.HandleFunc("/login", loginPageHandler.ServeHTTP)

	loginHandler := login.LoginHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config}
	http.HandleFunc("/user_login", loginHandler.ServeHTTP)

	logoutHandler := logout.LogoutHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config}
	logoutWrapper := wrappers.AuthenticationHandler{Next: logoutHandler, UserContext: handlerManager.UserContext, Config: handlerManager.Config}
	http.HandleFunc("/user_logout", logoutWrapper.ServeHTTP)

	ticketExplorerHandler := ticketexplorer.TicketExplorerPageHandler{UserContext: handlerManager.UserContext, TicketContext: handlerManager.TicketContext, Config: handlerManager.Config}
	http.HandleFunc("/tickets", ticketExplorerHandler.ServeHTTP)

	ticketViewHandler := ticketView.TicketViewPageHandler{UserContext: handlerManager.UserContext, TicketContext: handlerManager.TicketContext, Config: handlerManager.Config}
	http.HandleFunc("/ticket/", ticketViewHandler.ServeHTTP)
}
