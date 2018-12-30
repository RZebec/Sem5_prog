package wrappers

import (
	"context"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
)

const isAdminKey = "IsAdmin"
const isAuthenticatedKey = "IsAuthenticated"

func newContextWithAuthenticationInfo(ctx context.Context, isAuthenticated bool, isAdmin bool) context.Context {
	ctx = context.WithValue(ctx, isAdminKey, isAdmin)
	return context.WithValue(ctx, isAuthenticatedKey, isAuthenticated)
}

func isAdmin(ctx context.Context) bool {
	return ctx.Value(isAdminKey).(bool)
}

func isAuthenticated(ctx context.Context) bool {
	return ctx.Value(isAuthenticatedKey).(bool)
}

/*
	Function used to check if a user is logged in and if the session of the aforementioned user is valid.
*/
func UserIsLoggedInCheck(r *http.Request, userContext user.UserContext, accessTokenCookieName string, logger logging.Logger) (isUserLoggedIn bool, isAdmin bool, accessTokenValue string) {
	userIsLoggedIn := false
	token := ""

	cookie, err := r.Cookie(accessTokenCookieName)

	if err != nil {
		logger.LogError("Login", err)
		return false,  false,""
	}

	token = cookie.Value
	userIsLoggedIn, _, _, role, err := userContext.SessionIsValid(token)

	return userIsLoggedIn, role == user.Admin, token
}