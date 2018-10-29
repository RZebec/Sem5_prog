package helpers

import "net/http"

/*
	Sets the Cookie, with its respective name and value, in the browser.
 */
func SetCookie(w http.ResponseWriter, r *http.Request, name string, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

/*
	Removes the Cookie from the browser.
	It doesn't remove the cookie completely, just sets its value to a empty string.
 */
func RemoveCookie(w http.ResponseWriter, r *http.Request, name string) {
	SetCookie(w, r, name, "")
}
