package webui

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
