// 5894619, 6720876, 9793350
package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	Correct authentication info should be added.
*/
func TestEnforceAuthenticationInfoWrapper_ServeHTTP_UserIsAuthenticatedAndAdmin(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	nextHandler := testhelpers.LoggingHTPPHandler{}
	testee := EnforceAuthenticationWrapper{UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger(),
		Next:   &nextHandler}

	// Mock the context:
	mockedUserContext.On("RefreshToken", mock.Anything).Return("tete", nil)
	mockedUserContext.On("SessionIsValid", mock.Anything).Return(true, 2,
		"testName", userData.Admin, nil)

	req, err := http.NewRequest("GET", shared.ReceivePath, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", shared.AccessTokenCookieName+"=1234568")

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	handler.ServeHTTP(rr, req)

	// Assert that the next handler has been called and check the injected values of the context:
	assert.True(t, nextHandler.HasBeenCalled, "The next handler should be called")
	isAdmin := IsAdmin(nextHandler.Request.Context())
	isAuthenticated := IsAuthenticated(nextHandler.Request.Context())
	assert.True(t, isAdmin, "The next handler should get the info that the user is a admin")
	assert.True(t, isAuthenticated, "The next handler should get the info that the user is authenticated")

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")
	mockedUserContext.AssertExpectations(t)
}

/*
	Combination if user is authenticated but no admin should work.
*/
func TestEnforceAuthenticationInfoWrapper_ServeHTTP_UserIsAuthenticatedAndNoAdmin(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	nextHandler := testhelpers.LoggingHTPPHandler{}
	testee := EnforceAuthenticationWrapper{UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger(),
		Next:   &nextHandler}

	// Mock the context:
	mockedUserContext.On("RefreshToken", mock.Anything).Return("tete", nil)
	mockedUserContext.On("SessionIsValid", mock.Anything).Return(true, 2,
		"testName", userData.RegisteredUser, nil)

	req, err := http.NewRequest("GET", shared.ReceivePath, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", shared.AccessTokenCookieName+"=1234568")

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	handler.ServeHTTP(rr, req)

	// Assert that the next handler has been called and check the injected values of the context:
	assert.True(t, nextHandler.HasBeenCalled, "The next handler should be called")
	isAdmin := IsAdmin(nextHandler.Request.Context())
	isAuthenticated := IsAuthenticated(nextHandler.Request.Context())
	// The user is no admin:
	assert.False(t, isAdmin, "The next handler should get the info that the user is a no admin")
	assert.True(t, isAuthenticated, "The next handler should get the info that the user is authenticated")

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")
	mockedUserContext.AssertExpectations(t)
}

/*
	A user which is not authenticated.
*/
func TestEnforceAuthenticationInfoWrapper_ServeHTTP_NotAuthenticated(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	nextHandler := testhelpers.LoggingHTPPHandler{}
	testee := EnforceAuthenticationWrapper{UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger(),
		Next:   &nextHandler}

	// Mock the context:
	mockedUserContext.On("SessionIsValid", mock.Anything).Return(false, -1,
		"", userData.RegisteredUser, nil)

	req, err := http.NewRequest("GET", shared.ReceivePath, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", shared.AccessTokenCookieName+"=1234568")

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	handler.ServeHTTP(rr, req)

	// Assert that the next handler has not been called and check the injected values of the context:
	assert.False(t, nextHandler.HasBeenCalled, "The next handler should not be called")
	assert.Equal(t, 302, rr.Code, "Status code 302 should be returned")
	mockedUserContext.AssertExpectations(t)
}

/*
	A missing cookie should result in a non-authenticated user.
*/
func TestEnforceAuthenticationInfoWrapper_ServeHTTP_NotCookieSet(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	nextHandler := testhelpers.LoggingHTPPHandler{}
	testee := EnforceAuthenticationWrapper{UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger(),
		Next:   &nextHandler}

	req, err := http.NewRequest("GET", shared.ReceivePath, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	handler.ServeHTTP(rr, req)

	// Assert that the next handler has been called and check the injected values of the context:
	assert.False(t, nextHandler.HasBeenCalled, "The next handler should not be called")

	assert.Equal(t, 302, rr.Code, "Status code 302 should be returned")
	mockedUserContext.AssertExpectations(t)
}
