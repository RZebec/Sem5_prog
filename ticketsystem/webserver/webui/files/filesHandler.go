package files

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files/script"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files/style"
	"net/http"
	"strings"
)

type FilesHandler struct {
}

func (f FilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")

	switch s[2] {
	case "style":
		style.HandelStyle(w, r)
	case "script":
		script.HandelScript(w, r)
	}
}
