package helpers

import "net/http"

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
