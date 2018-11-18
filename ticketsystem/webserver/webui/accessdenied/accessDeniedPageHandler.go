package accessdenied

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"fmt"
	"net/http"
)

/*
	Structure for the Access Denied Page Handler.
*/
type AccessDeniedPageHandler struct {
}

/*
	Structure for the Access Denied Page Data.
*/
type accessDeniedPageData struct {
}

/*
	The Access Denied Page handler.
*/
func (l AccessDeniedPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data := accessDeniedPageData{}

	err := templateManager.RenderTemplate(w, "AccessDeniedPage", data)

	if err != nil {
		// TODO: Handle error
		fmt.Print(err)
	}
}
