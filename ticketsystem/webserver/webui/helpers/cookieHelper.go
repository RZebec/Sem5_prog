package helpers

import (
	"net/http"
)

type Cookie struct {
	Name  string
	Value string
}

func (c Cookie) SetCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func (c Cookie) RemoveCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     c.Name,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func (c Cookie) RefreshCookieValue(r *http.Request) {
	cookie, err := r.Cookie(c.Name)

	if err != nil {
		// TODO: Handle error
	} else {
		c.Value = cookie.Value
	}
}
