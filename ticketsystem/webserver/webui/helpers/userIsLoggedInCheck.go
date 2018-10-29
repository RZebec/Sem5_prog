package helpers

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/session"
	"net/http"
)

/*
	Function used to check if a user is logged in and if the session of the aforementioned user is valid.
 */
func UserIsLoggedInCheck(r *http.Request, sessionManager session.SessionManager, accessTokenCookie Cookie) (isUserLoggedIn bool, accessTokenValue string) {
	userIsLoggedIn := false
	token := ""

	cookie, err := r.Cookie(accessTokenCookie.Name)

	if err == nil {
		token = cookie.Value
		userIsLoggedIn, _, _, err = sessionManager.SessionIsValid(token)
	}

	return userIsLoggedIn, token
}
