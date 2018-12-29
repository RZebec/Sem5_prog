package webui

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/admin"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/login"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/logout"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/register"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
)

type HandlerManager struct {
	UserContext   user.UserContext
	TicketContext ticket.TicketContext
	Config        config.Configuration
	Logger        logging.Logger
}

func (handlerManager *HandlerManager) StartServices() {

	filesHandler := files.FilesHandler{}
	http.HandleFunc("/files/", filesHandler.ServeHTTP)

	registerHandler := register.RegisterHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/register", registerHandler.ServeHTTPGetRegisterPage)
	http.HandleFunc("/user_register", registerHandler.ServeHTTPPostRegisteringData)

	loginHandler := login.LoginHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/login", loginHandler.ServeHTTPGetLoginPage)
	http.HandleFunc("/user_login", loginHandler.ServeHTTPPostLoginData)

	logoutHandler := logout.LogoutHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	logoutWrapper := wrappers.AuthenticationWrapper{Next: logoutHandler, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/user_logout", logoutWrapper.ServeHTTP)

	adminPageHandler := admin.AdminPageHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	adminPageWrapper := wrappers.AdminWrapper{Next: adminPageHandler, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	adminPageAuthenticationWrapper := wrappers.AuthenticationWrapper{Next: adminPageWrapper, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/admin", adminPageAuthenticationWrapper.ServeHTTP)

	adminSetApiKeysHandler := admin.AdminSetApiKeysHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	adminSetApiKeysWrapper := wrappers.AdminWrapper{Next: adminSetApiKeysHandler, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	adminSetApiKeysAuthenticationWrapper := wrappers.AuthenticationWrapper{Next: adminSetApiKeysWrapper, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/set_api_keys", adminSetApiKeysAuthenticationWrapper.ServeHTTP)

	adminUnlockUserHandler := admin.AdminUnlockUserHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	adminUnlockUserWrapper := wrappers.AdminWrapper{Next: adminUnlockUserHandler, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	adminUnlockUserAuthenticationWrapper := wrappers.AuthenticationWrapper{Next: adminUnlockUserWrapper, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/unlock_user/", adminUnlockUserAuthenticationWrapper.ServeHTTP)

}
