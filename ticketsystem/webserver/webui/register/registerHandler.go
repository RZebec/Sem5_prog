package register

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"html"
	"net/http"
)

/*
	Structure for the User Register handler.
*/
type UserRegisterHandler struct {
	UserContext     userData.UserContext
	Logger          logging.Logger
}

/*
	The User Register handler.
*/
func (l UserRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			l.Logger.LogInfo("UserRegisterHandler","User registered" )
		} else {
			http.Redirect(w, r, "/register?IsRegisteringFailed=true", http.StatusSeeOther)
			l.Logger.LogInfo("UserRegisterHandler","User registration failed" )
		}
	}
}
