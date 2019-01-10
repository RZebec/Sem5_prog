// 5894619, 6720876, 9793350
package logout

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	Logout of a logged in user should be possible.
*/
func TestLogoutHandler_ServeHTTP_UserWasLoggedIn(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	cookieValue := "etete3tas"
	testee := UserLogoutHandler{UserContext: mockedUserContext, Logger: testhelpers.GetTestLogger()}
	mockedUserContext.On("Logout", cookieValue)

	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", shared.AccessTokenCookieName+"="+cookieValue)
	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	handler.ServeHTTP(rr, req)

	// Assert that the cookie value is empty:
	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusFound, resp.StatusCode, "Should return status code 302")
	assert.Equal(t, "/", newLocation, "Should be redirected to /")

	cookieExists, cookieValue := testhelpers.GetCookieValue(resp.Cookies(), shared.AccessTokenCookieName)
	assert.True(t, cookieExists, "The cookie should be set")
	assert.Equal(t, "", cookieValue, "The cookie should be empty")

	mockedUserContext.AssertExpectations(t)
}

/*
	Logout without a logged in user should be possible.
*/
func TestLogoutHandler_ServeHTTP_UserWasNotLoggedIn(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	testee := UserLogoutHandler{UserContext: mockedUserContext, Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	handler.ServeHTTP(rr, req)

	// Assert that the cookie value is empty:
	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusFound, resp.StatusCode, "Should return status code 302")
	assert.Equal(t, "/", newLocation, "Should be redirected to /")

	cookieExists, cookieValue := testhelpers.GetCookieValue(resp.Cookies(), shared.AccessTokenCookieName)
	assert.True(t, cookieExists, "The cookie should be set")
	assert.Equal(t, "", cookieValue, "The cookie should be empty")

	mockedUserContext.AssertExpectations(t)
}
