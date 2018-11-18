package register

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"fmt"
	"net/http"
	"strconv"
)

/*
	Structure for the Register Page Handler.
*/
type RegisterPageHandler struct {
	Config      config.Configuration
	UserContext user.UserContext
}

/*
	Structure for the Register Page Data.
*/
type registerPageData struct {
	IsRegisteringFailed bool
}

/*
	The Register Page handler.
*/
func (l RegisterPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Checks if the User is already logged in and if so redirects him to the start page
	isUserLoggedIn, _ := helpers.UserIsLoggedInCheck(r, l.UserContext, l.Config.AccessTokenCookieName)

	if isUserLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	queryValues := r.URL.Query()
	isRegisteringFailed, err := strconv.ParseBool(queryValues.Get("IsRegisteringFailed"))

	if err != nil {
		// TODO: Handle error
	}

	data := registerPageData{
		IsRegisteringFailed: isRegisteringFailed,
	}

	err = templateManager.RenderTemplate(w, "RegisterPage", data)

	if err != nil {
		// TODO: Handle error
		fmt.Print(err)
	}
}
