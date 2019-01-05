package wrappers

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	A request for a admin should be relayed to the next handler.
*/
func TestAdminWrapper_ServeHTTP_IsAdmin(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	nextHandler := testhelpers.LoggingHTPPHandler{}
	testee := AdminWrapper{UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger(),
		Next:   &nextHandler}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContextWithAuthenticationInfo(req.Context(), true, true, 1, "")

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	// Assert that the next handler has been called:
	assert.True(t, nextHandler.HasBeenCalled, "The next handler should be called")
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")
	mockedUserContext.AssertExpectations(t)
}

/*
	A request for non-admin should not be relayed.
*/
func TestAdminWrapper_ServeHTTP_IsNoAdmin(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	nextHandler := testhelpers.LoggingHTPPHandler{}
	testee := AdminWrapper{UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger(),
		Next:   &nextHandler}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Should be no admin:
	ctx := NewContextWithAuthenticationInfo(req.Context(), true, false, 1, "")

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	// Assert that the next handler has not been called:
	assert.False(t, nextHandler.HasBeenCalled, "The next handler should not be called")
	assert.Equal(t, 403, rr.Code, "Status code 403 should be returned")
	mockedUserContext.AssertExpectations(t)
}

/*
	A request without authentication info should not be relayed to the next handler.
*/
func TestAdminWrapper_ServeHTTP_NoValuesInContextSet(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	nextHandler := testhelpers.LoggingHTPPHandler{}
	testee := AdminWrapper{UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger(),
		Next:   &nextHandler}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test with no authentication info set:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	handler.ServeHTTP(rr, req)

	// Assert that the next handler has not been called:
	assert.False(t, nextHandler.HasBeenCalled, "The next handler should not be called")
	assert.Equal(t, 403, rr.Code, "Status code 403 should be returned")
	mockedUserContext.AssertExpectations(t)
}
