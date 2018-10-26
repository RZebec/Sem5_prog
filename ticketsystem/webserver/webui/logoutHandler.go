package webui

import (
	"fmt"
	"net/http"
	"strings"
)

type LogoutHandler struct {
}

func (l LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		fmt.Println("User Logged out")
	}
}
