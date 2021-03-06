// 5894619, 6720876, 9793350
package files

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files/script"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/files/style"
	"net/http"
	"strings"
)

/*
	The struct for the files handler.
*/
type FileHandler struct {
}

/*
	The handler for the different files.
	Differs between the style(CSS) and the script(Javascript) files.
*/
func (f FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")

	switch s[2] {
	case "style":
		style.HandelStyle(w, r)
	case "script":
		script.HandelScript(w, r)
	}
}
