package webui

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/admin"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/index"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/login"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/logout"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/register"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
)

type HandlerManager struct {
	UserContext   user.UserContext
	TicketContext ticket.TicketContext
	Config        config.Configuration
	Logger        logging.Logger
	ApiConfiguration	config.ApiContext
	TemplateManager		templateManager.TemplateContext
}

func (handlerManager *HandlerManager) RegisterHandlers() {

	filesHandler := files.FilesHandler{}
	http.HandleFunc("/files/", filesHandler.ServeHTTP)

	indexPageHandler := index.IndexPageHandler{Logger: handlerManager.Logger, TemplateManager: handlerManager.TemplateManager}
	indexPageAuthenticationInfoWrapper := wrappers.AddAuthenticationInfoWrapper{}
	indexPageAuthenticationInfoWrapper.Next = indexPageHandler
	indexPageAuthenticationInfoWrapper.Logger = handlerManager.Logger
	indexPageAuthenticationInfoWrapper.UserContext = handlerManager.UserContext

	http.HandleFunc("/", indexPageAuthenticationInfoWrapper.ServeHTTP)

	registerHandler := register.RegisterHandler{UserContext: handlerManager.UserContext, Logger: handlerManager.Logger, TemplateManager: handlerManager.TemplateManager}
	http.HandleFunc("/register", registerHandler.ServeHTTPGetRegisterPage)
	http.HandleFunc("/user_register", registerHandler.ServeHTTPPostRegisteringData)

	loginHandler := login.LoginHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger, TemplateManager: handlerManager.TemplateManager}
	http.HandleFunc("/login", loginHandler.ServeHTTPGetLoginPage)
	http.HandleFunc("/user_login", loginHandler.ServeHTTPPostLoginData)

	logoutHandler := logout.LogoutHandler{UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	logoutWrapper := wrappers.EnforceAuthenticationWrapper{Next: logoutHandler, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/user_logout", logoutWrapper.ServeHTTP)

	adminPageHandler := admin.AdminPageHandler{UserContext: handlerManager.UserContext, Logger: handlerManager.Logger, TemplateManager: handlerManager.TemplateManager, ApiContext: handlerManager.ApiConfiguration}
	adminPageWrapper := wrappers.AdminWrapper{Next: adminPageHandler, UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	adminPageAuthenticationWrapper := wrappers.EnforceAuthenticationWrapper{Next: adminPageWrapper, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/admin", adminPageAuthenticationWrapper.ServeHTTP)

	adminSetApiKeysHandler := admin.AdminSetApiKeysHandler{Logger: handlerManager.Logger, ApiConfiguration: handlerManager.ApiConfiguration}
	adminSetApiKeysWrapper := wrappers.AdminWrapper{Next: adminSetApiKeysHandler, UserContext: handlerManager.UserContext,  Logger: handlerManager.Logger}
	adminSetApiKeysAuthenticationWrapper := wrappers.EnforceAuthenticationWrapper{Next: adminSetApiKeysWrapper, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/set_api_keys", adminSetApiKeysAuthenticationWrapper.ServeHTTP)

	adminUnlockUserHandler := admin.AdminUnlockUserHandler{UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	adminUnlockUserWrapper := wrappers.AdminWrapper{Next: adminUnlockUserHandler, UserContext: handlerManager.UserContext,  Logger: handlerManager.Logger}
	adminUnlockUserAuthenticationWrapper := wrappers.EnforceAuthenticationWrapper{Next: adminUnlockUserWrapper, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/unlock_user", adminUnlockUserAuthenticationWrapper.ServeHTTP)

}
