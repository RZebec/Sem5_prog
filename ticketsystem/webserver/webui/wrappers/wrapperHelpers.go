package wrappers

import (
	"context"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"net/http"
)

const isAdminKey = "IsAdmin"
const isAuthenticatedKey = "IsAuthenticated"
const userIdKey = "userId"

/*
	Inject the context with authentication info.
*/
func NewContextWithAuthenticationInfo(ctx context.Context, isAuthenticated bool, isAdmin bool, userId int) context.Context {
	ctx = context.WithValue(ctx, isAdminKey, isAdmin)
	ctx = context.WithValue(ctx, isAuthenticatedKey, isAuthenticated)
	return context.WithValue(ctx, userIdKey, userId)
}

/*
	Return true, if the user is a admin.
*/
func IsAdmin(ctx context.Context) bool {
	value, ok := ctx.Value(isAdminKey).(bool)
	if ok {
		return value
	} else {
		return false
	}
}

/*
	Returns true if the user is authenticated.
*/
func IsAuthenticated(ctx context.Context) bool {
	value, ok := ctx.Value(isAuthenticatedKey).(bool)
	if ok {
		return value
	} else {
		return false
	}
}

/*
	Get the user Id.
 */
func GetUserId(ctx context.Context) int {
	value, ok := ctx.Value(userIdKey).(int)
	if ok {
		return value
	} else {
		return -1
	}
}

/*
	Function used to check if a user is logged in and if the session of the aforementioned user is valid.
*/
func UserIsLoggedInCheck(r *http.Request, userContext user.UserContext, accessTokenCookieName string, logger logging.Logger) (isUserLoggedIn bool, isAdmin bool, accessTokenValue string, userId int) {
	userIsLoggedIn := false
	token := ""

	cookie, err := r.Cookie(accessTokenCookieName)

	if err != nil {
		logger.LogError("UserIsLoggedInCheck", err)
		return false, false, "", -1
	}

	token = cookie.Value
	userIsLoggedIn, userId, _, role, err := userContext.SessionIsValid(token)

	return userIsLoggedIn, role == user.Admin, token, userId
}
