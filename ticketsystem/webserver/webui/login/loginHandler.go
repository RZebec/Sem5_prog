package login

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui"
	"fmt"
	"net/http"
	"strings"
)

type LoginHandler struct {
	UserManager      session.UserManager
	LoginPageHandler LoginPageHandler
}

func (l LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		userName := r.FormValue("userName")
		password := r.FormValue("password")

		//success, err := l.UserManager.Register(userName, password)

		//if err != nil {
		//	fmt.Println(err.Error())
		//}

		success, token, err := l.UserManager.Login(userName, password)

		if err != nil {
			fmt.Println(err.Error())
		}

		if success {
			c := webui.Cookie{Name: "Access-Token", Value: token}
			c.SetCookie(w, r)
			http.Redirect(w, r, "/", 302)
		} else {
			http.Redirect(w, r, "/login?IsLoginFailed=true", 302)
		}
	}
}
