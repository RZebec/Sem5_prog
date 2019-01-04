package webui

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/admin"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/index"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/login"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/logout"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/register"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/tickets"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/userSettings"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
)

type HandlerManager struct {
	UserContext      user.UserContext
	TicketContext    ticket.TicketContext
	Config           config.Configuration
	Logger           logging.Logger
	ApiConfiguration config.ApiContext
	TemplateManager  templateManager.TemplateContext
	MailContext 	 mail.MailContext
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

	loginHandler := login.LoginHandler{UserContext: handlerManager.UserContext, Logger: handlerManager.Logger, TemplateManager: handlerManager.TemplateManager}
	http.HandleFunc("/login", loginHandler.ServeHTTPGetLoginPage)
	http.HandleFunc("/user_login", loginHandler.ServeHTTPPostLoginData)

	logoutHandler := logout.LogoutHandler{UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	logoutWrapper := wrappers.EnforceAuthenticationWrapper{Next: logoutHandler, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/user_logout", logoutWrapper.ServeHTTP)

	adminPageHandler := admin.AdminPageHandler{UserContext: handlerManager.UserContext, Logger: handlerManager.Logger, TemplateManager: handlerManager.TemplateManager, ApiContext: handlerManager.ApiConfiguration}
	adminPageWrapper := wrappers.AdminWrapper{Next: adminPageHandler, UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	adminPageAuthenticationWrapper := wrappers.EnforceAuthenticationWrapper{Next: adminPageWrapper, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/admin", adminPageAuthenticationWrapper.ServeHTTP)

	adminSetApiKeysHandler := admin.AdminSetApiKeysHandler{Logger: handlerManager.Logger, ApiConfiguration: handlerManager.ApiConfiguration}
	adminSetApiKeysWrapper := wrappers.AdminWrapper{Next: adminSetApiKeysHandler, UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	adminSetApiKeysAuthenticationWrapper := wrappers.EnforceAuthenticationWrapper{Next: adminSetApiKeysWrapper, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/set_api_keys", adminSetApiKeysAuthenticationWrapper.ServeHTTP)

	adminUnlockUserHandler := admin.AdminUnlockUserHandler{UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	adminUnlockUserWrapper := wrappers.AdminWrapper{Next: adminUnlockUserHandler, UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	adminUnlockUserAuthenticationWrapper := wrappers.EnforceAuthenticationWrapper{Next: adminUnlockUserWrapper, UserContext: handlerManager.UserContext, Config: handlerManager.Config, Logger: handlerManager.Logger}
	http.HandleFunc("/unlock_user", adminUnlockUserAuthenticationWrapper.ServeHTTP)

	allTicketExplorerPageHandler := tickets.AllTicketsExplorerPageHandler{TicketContext: handlerManager.TicketContext, TemplateManager: handlerManager.TemplateManager, Logger: handlerManager.Logger}
	allTicketExplorerPageHandlerWrapper := wrappers.AddAuthenticationInfoWrapper{}
	allTicketExplorerPageHandlerWrapper.Next = allTicketExplorerPageHandler
	allTicketExplorerPageHandlerWrapper.Logger = handlerManager.Logger
	allTicketExplorerPageHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/all_tickets", allTicketExplorerPageHandlerWrapper.ServeHTTP)

	openTicketExplorerPageHandler := tickets.OpenTicketsExplorerPageHandler{TicketContext: handlerManager.TicketContext, TemplateManager: handlerManager.TemplateManager, Logger: handlerManager.Logger}
	openTicketExplorerPageHandlerWrapper := wrappers.AddAuthenticationInfoWrapper{}
	openTicketExplorerPageHandlerWrapper.Next = openTicketExplorerPageHandler
	openTicketExplorerPageHandlerWrapper.Logger = handlerManager.Logger
	openTicketExplorerPageHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/open_tickets", openTicketExplorerPageHandlerWrapper.ServeHTTP)

	userTicketExplorerPageHandler := tickets.UserTicketsExplorerPageHandler{TicketContext: handlerManager.TicketContext, TemplateManager: handlerManager.TemplateManager, Logger: handlerManager.Logger}
	userTicketExplorerPageHandlerWrapper := wrappers.AddAuthenticationInfoWrapper{}
	userTicketExplorerPageHandlerWrapper.Next = userTicketExplorerPageHandler
	userTicketExplorerPageHandlerWrapper.Logger = handlerManager.Logger
	userTicketExplorerPageHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/user_tickets", userTicketExplorerPageHandlerWrapper.ServeHTTP)

	ticketViewPageHandler := tickets.TicketViewPageHandler{UserContext: handlerManager.UserContext, TicketContext: handlerManager.TicketContext, TemplateManager: handlerManager.TemplateManager, Logger: handlerManager.Logger}
	ticketViewPageHandlerWrapper := wrappers.AddAuthenticationInfoWrapper{}
	ticketViewPageHandlerWrapper.Next = ticketViewPageHandler
	ticketViewPageHandlerWrapper.Logger = handlerManager.Logger
	ticketViewPageHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/ticket/", ticketViewPageHandlerWrapper.ServeHTTP)

	ticketAppendMessageHandler := tickets.TicketAppendMessageHandler{TicketContext: handlerManager.TicketContext,
		UserContext: handlerManager.UserContext, MailContext: handlerManager.MailContext, Logger: handlerManager.Logger}
	ticketAppendHandlerWrapper := wrappers.AddAuthenticationInfoWrapper{}
	ticketAppendHandlerWrapper.UserContext = handlerManager.UserContext
	ticketAppendHandlerWrapper.Logger = handlerManager.Logger
	ticketAppendHandlerWrapper.Next = ticketAppendMessageHandler
	http.HandleFunc("/append_message", ticketAppendHandlerWrapper.ServeHTTP)

	ticketMergeHandler := tickets.TicketMergeHandler{TicketContext: handlerManager.TicketContext,
		MailContext: handlerManager.MailContext, Logger: handlerManager.Logger}
	ticketMergeEnforceAuthenticationWrapper:= wrappers.EnforceAuthenticationWrapper{}
	ticketMergeEnforceAuthenticationWrapper.Next = ticketMergeHandler
	ticketMergeEnforceAuthenticationWrapper.Logger = handlerManager.Logger
	ticketMergeEnforceAuthenticationWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/merge_tickets", ticketMergeEnforceAuthenticationWrapper.ServeHTTP)

	ticketSetEditorHandler := tickets.TicketSetEditorHandler{TicketContext: handlerManager.TicketContext,
		MailContext: handlerManager.MailContext, UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	setEditorWrapper := wrappers.EnforceAuthenticationWrapper{}
	setEditorWrapper.Next = ticketSetEditorHandler
	setEditorWrapper.Logger = handlerManager.Logger
	setEditorWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/ticket_setEditor", setEditorWrapper.ServeHTTP)

	userSettingsPageHandler := userSettings.UserSettingsPageHandler{UserContext: handlerManager.UserContext, TemplateManager: handlerManager.TemplateManager, Logger: handlerManager.Logger}
	userSettingsPageHandlerWrapper := wrappers.EnforceAuthenticationWrapper{}
	userSettingsPageHandlerWrapper.Next = userSettingsPageHandler
	userSettingsPageHandlerWrapper.Logger = handlerManager.Logger
	userSettingsPageHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/user_settings", userSettingsPageHandlerWrapper.ServeHTTP)

	changePasswordHandler := userSettings.ChangePasswordHandler{UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	changePasswordHandlerWrapper := wrappers.EnforceAuthenticationWrapper{}
	changePasswordHandlerWrapper.Next = changePasswordHandler
	changePasswordHandlerWrapper.Logger = handlerManager.Logger
	changePasswordHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/user_change_password", changePasswordHandlerWrapper.ServeHTTP)

	toggleVacationModeHandler := userSettings.ToggleVacationModeHandler{UserContext: handlerManager.UserContext, Logger: handlerManager.Logger}
	toggleVacationModeHandlerWrapper := wrappers.EnforceAuthenticationWrapper{}
	toggleVacationModeHandlerWrapper.Next = toggleVacationModeHandler
	toggleVacationModeHandlerWrapper.Logger = handlerManager.Logger
	toggleVacationModeHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/user_toggle_vacation", toggleVacationModeHandlerWrapper.ServeHTTP)

	ticketCreatePageHandler := tickets.TicketCreatePageHandler{UserContext: handlerManager.UserContext, TemplateManager: handlerManager.TemplateManager, Logger: handlerManager.Logger}
	ticketCreatePageHandlerWrapper := wrappers.AddAuthenticationInfoWrapper{}
	ticketCreatePageHandlerWrapper.Next = ticketCreatePageHandler
	ticketCreatePageHandlerWrapper.Logger = handlerManager.Logger
	ticketCreatePageHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/ticket_create", ticketCreatePageHandlerWrapper.ServeHTTP)

	ticketCreateHandler := tickets.TicketCreateHandler{UserContext: handlerManager.UserContext, TicketContext: handlerManager.TicketContext, Logger: handlerManager.Logger}
	ticketCreateHandlerWrapper := wrappers.AddAuthenticationInfoWrapper{}
	ticketCreateHandlerWrapper.Next = ticketCreateHandler
	ticketCreateHandlerWrapper.Logger = handlerManager.Logger
	ticketCreateHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/create_ticket", ticketCreateHandlerWrapper.ServeHTTP)

	ticketEditPageHandler := tickets.TicketEditPageHandler{UserContext: handlerManager.UserContext, TicketContext: handlerManager.TicketContext, Logger: handlerManager.Logger, TemplateManager: handlerManager.TemplateManager}
	ticketEditPageHandlerWrapper := wrappers.EnforceAuthenticationWrapper{}
	ticketEditPageHandlerWrapper.Next = ticketEditPageHandler
	ticketEditPageHandlerWrapper.Logger = handlerManager.Logger
	ticketEditPageHandlerWrapper.UserContext = handlerManager.UserContext
	http.HandleFunc("/ticket/ticket_edit/", ticketEditPageHandlerWrapper.ServeHTTP)
}
