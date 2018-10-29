package helpers

import "net/http"

/*
	Structure for the Cookie.
 */
type Cookie struct {
	Name  string
	Value string
}

/*
	Sets the Cookie, with its respective name and value, in the browser.
 */
func (c Cookie) SetCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

/*
	Removes the Cookie from the browser.
	It doesn't remove the cookie completely, just sets its value to a empty string.
 */
func (c Cookie) RemoveCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     c.Name,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}
