package register

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"html"
	"net/http"
	"strconv"
)

/*
	Structure for the Register handler.
*/
type RegisterHandler struct {
	UserContext     userData.UserContext
	Logger          logging.Logger
	TemplateManager templateManager.TemplateContext
}

/*
	Structure for the Register Page Data.
*/
type registerPageData struct {
	IsRegisteringFailed bool
	pages.BasePageData
}

/*
	The Login handler.
*/
func (l RegisterHandler) ServeHTTPPostRegisteringData(w http.ResponseWriter, r *http.Request) {
	// TODO: Verification step for the userData needs to be implemented here
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		userName := r.FormValue("userName")
		password := r.FormValue("password")

		firstName = html.EscapeString(firstName)
		lastName = html.EscapeString(lastName)
		userName = html.EscapeString(userName)
		password = html.EscapeString(password)

		success, err := l.UserContext.Register(userName, password, firstName, lastName)

		if err != nil {
			l.Logger.LogError("Register", err)
			http.Redirect(w, r, "/register?IsRegisteringFailed=true", http.StatusSeeOther)
		}

		if success {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			l.Logger.LogInfo("RegisterHandler","User registered" )
		} else {
			http.Redirect(w, r, "/register?IsRegisteringFailed=true", http.StatusSeeOther)
			l.Logger.LogInfo("RegisterHandler","User registration failed" )
		}
	}
}

/*
	The Register Page handler.
*/
func (l RegisterHandler) ServeHTTPGetRegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		isUserLoggedIn := wrappers.IsAuthenticated(r.Context())
		userIsAdmin := wrappers.IsAdmin(r.Context())

		if isUserLoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		queryValues := r.URL.Query()
		queryValue := queryValues.Get("IsRegisteringFailed")
		isRegisteringFailed, err := strconv.ParseBool(queryValue)

		if err != nil && queryValue != "" {
			l.Logger.LogError("Register", err)
			isRegisteringFailed = false
		}

		data := registerPageData{
			IsRegisteringFailed: isRegisteringFailed,
		}

		data.UserIsAuthenticated = wrappers.IsAuthenticated(r.Context())
		data.UserIsAdmin = userIsAdmin
		data.Active = "register"

		err = l.TemplateManager.RenderTemplate(w, "RegisterPage", data)

		if err != nil {
			l.Logger.LogError("Register", err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
	}
}
