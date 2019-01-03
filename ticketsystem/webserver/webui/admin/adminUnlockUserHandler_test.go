package admin

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

/*
	Only POST request should be possible.
*/
func TestAdminUnlockUserHandlerWrongRequestMethod_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("GET", "/set_api_keys", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := AdminUnlockUserHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Status code 405 should be returned")
	mockedUserContext.AssertExpectations(t)
}

/*
	A request with incorrect data should return a 400.
*/
func TestAdminUnlockUserHandle_ServeHTTP_IncorrectData(t *testing.T) {
	req, err := http.NewRequest("POST", "/set_api_keys", nil)
	req.Form = url.Values{}
	req.Form.Add("userId", "asdhajs")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := AdminUnlockUserHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code 400 should be returned")
	mockedUserContext.AssertExpectations(t)
}

/*
	A request with incorrect cookie should return a 400
*/
func TestAdminUnlockUserHandle_ServeHTTP_NoPresentCookie(t *testing.T) {
	req, err := http.NewRequest("POST", "/set_api_keys", nil)
	req.Form = url.Values{}
	req.Form.Add("userId", "1234")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := AdminUnlockUserHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code 400 should be returned")

	mockedUserContext.AssertExpectations(t)
}

/*
	A valid request should be possible.
*/
func TestAdminUnlockUserHandle_ServeHTTP_ValidRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "/set_api_keys", nil)
	req.Form = url.Values{}
	req.Form.Add("userId", "1234")
	if err != nil {
		t.Fatal(err)
	}

	cookie := http.Cookie{Name: shared.AccessTokenCookieName, Value: "test"}
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("UnlockAccount", mock.Anything, mock.Anything).Return(true, nil)

	testee := AdminUnlockUserHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedUserContext.AssertExpectations(t)
}

/*
	A error from changing the outgoing api key should result in a 500.
*/
func TestAdminUnlockUserHandle_ServeHTTP_UnlockAccount_ContextReturnError_500Returned(t *testing.T) {
	req, err := http.NewRequest("POST", "/set_api_keys", nil)
	req.Form = url.Values{}
	req.Form.Add("userId", "1234")
	if err != nil {
		t.Fatal(err)
	}

	cookie := http.Cookie{Name: shared.AccessTokenCookieName, Value: "test"}
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("UnlockAccount", mock.Anything, mock.Anything).Return(false, errors.New("TestError"))

	testee := AdminUnlockUserHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")

	mockedUserContext.AssertExpectations(t)
}

