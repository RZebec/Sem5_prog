package register

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/config"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
	"strings"
)

/*
	Structure for the Register handler.
*/
type RegisterHandler struct {
	UserContext user.UserContext
	Config      config.Configuration
}

/*
	The Login handler.
*/
func (reg RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Verification step for the user needs to be implemented here
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		userName := r.FormValue("userName")
		password := r.FormValue("password")

		success, err := reg.UserContext.Register(userName, password, firstName, lastName)

		if err != nil {
			// TODO: handle error
		}

		if success {
			http.Redirect(w, r, "/login", 302)
		} else {
			http.Redirect(w, r, "/register?IsRegisteringFailed=true", 302)
		}
	}
}
