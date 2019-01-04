package helpers

import (
	"net/http"
	"time"
)

/*
	Sets the Cookie, with its respective name and value, in the browser.
*/
func SetCookie(w http.ResponseWriter, name string, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

/*
	Remove a cookie.
 */
func RemoveCookie(w http.ResponseWriter, name string){
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires: time.Now().AddDate(0,-1,0),
	}
	http.SetCookie(w, &cookie)
}