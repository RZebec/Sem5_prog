package helpers

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
)

/*
	Function used to check if a user is logged in and if the session of the aforementioned user is valid.
 */
func UserIsLoggedInCheck(r *http.Request, userContext user.UserContext, accessTokenCookieName string) (isUserLoggedIn bool, accessTokenValue string) {
	userIsLoggedIn := false
	token := ""

	cookie, err := r.Cookie(accessTokenCookieName)

	if err == nil {
		token = cookie.Value
		userIsLoggedIn, _, _, err = userContext.SessionIsValid(token)
	}

	return userIsLoggedIn, token
}
