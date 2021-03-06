// 5894619, 6720876, 9793350
package wrappers

import (
	"context"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"net/http"
)

const isAdminKey = "IsAdmin"
const isAuthenticatedKey = "IsAuthenticated"
const userIdKey = "userId"
const userTokenKey = "userToken"

/*
	Inject the context with authentication info.
*/
func NewContextWithAuthenticationInfo(ctx context.Context, isAuthenticated bool, isAdmin bool, userId int, currentToken string) context.Context {
	ctx = context.WithValue(ctx, isAdminKey, isAdmin)
	ctx = context.WithValue(ctx, isAuthenticatedKey, isAuthenticated)
	ctx = context.WithValue(ctx, userTokenKey, currentToken)
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
	Get the user token.
*/
func GetUserToken(ctx context.Context) string {
	value, ok := ctx.Value(userTokenKey).(string)
	if ok {
		return value
	} else {
		return ""
	}
}

/*
	Function used to check if a user is logged in and if the session of the aforementioned user is valid.
*/
func UserIsLoggedInCheck(r *http.Request, userContext userData.UserContext, accessTokenCookieName string, logger logging.Logger) (isUserLoggedIn bool, isAdmin bool, accessTokenValue string, userId int) {
	userIsLoggedIn := false
	token := ""

	cookie, err := r.Cookie(accessTokenCookieName)

	if err != nil {
		if err.Error() == "http: named cookie not present" {
			return false, false, "", -1
		}
		logger.LogError("UserIsLoggedInCheck", err)
		return false, false, "", -1
	}

	token = cookie.Value
	userIsLoggedIn, userId, _, role, err := userContext.SessionIsValid(token)

	return userIsLoggedIn, role == userData.Admin, token, userId
}
