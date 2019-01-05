package register

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"net/http"
	"strconv"
)

/*
	Structure for the Register Page handler.
*/
type PageHandler struct {
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
	The Register Page handler.
*/
func (l PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
