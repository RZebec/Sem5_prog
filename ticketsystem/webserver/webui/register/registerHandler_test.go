// 5894619, 6720876, 9793350
package register

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

/*
	A valid request should be possible.
*/
func TestUserRegisterHandler_ServeHTTP_ValidRequest_RedirectedToLoginPage(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_register", nil)
	req.Form = url.Values{}
	req.Form.Add("first_name", "Max")
	req.Form.Add("last_name", "Muller")
	req.Form.Add("userName", "Max.Muller@test.de")
	req.Form.Add("password", "1234")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("Register", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

	testee := UserRegisterHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 303, rr.Code, "Status code 303 should be returned")
	assert.Equal(t, "/login", rr.Header().Get("location"), "User should be redirected to url \"/login\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	Only POST request should be possible.
*/
func TestUserRegisterHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/user_register", nil)
	req.Form = url.Values{}
	req.Form.Add("first_name", "Max")
	req.Form.Add("last_name", "Muller")
	req.Form.Add("userName", "Max.Muller@test.de")
	req.Form.Add("password", "1234")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := UserRegisterHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Status code 405 should be returned")

	mockedUserContext.AssertExpectations(t)
}

/*
	An unsuccessful registering procedure should return the same page.
*/
func TestUserRegisterHandler_ServeHTTP_UnsuccessfulRegistering_RedirectedToSamePageWithMessage(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_register", nil)
	req.Form = url.Values{}
	req.Form.Add("first_name", "Max")
	req.Form.Add("last_name", "Muller")
	req.Form.Add("userName", "Max.Muller@test.de")
	req.Form.Add("password", "1234")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("Register", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, nil)

	testee := UserRegisterHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	rr.Result()

	assert.Equal(t, 303, rr.Code, "Status code 303 should be returned")
	assert.Equal(t, "/register?IsRegisteringFailed=true", rr.Header().Get("location"), "User should be redirected to url \"/register?IsRegisteringFailed=true\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	An unsuccessful registering procedure should return the same page.
*/
func TestUserRegisterHandler_ServeHTTP_ContextReturnError_RedirectedToSamePageWithMessage(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_register", nil)
	req.Form = url.Values{}
	req.Form.Add("first_name", "Max")
	req.Form.Add("last_name", "Muller")
	req.Form.Add("userName", "Max.Muller@test.de")
	req.Form.Add("password", "1234")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("Register", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, errors.New("TestError"))

	testee := UserRegisterHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	rr.Result()

	assert.Equal(t, 303, rr.Code, "Status code 303 should be returned")
	assert.Equal(t, "/register?IsRegisteringFailed=true", rr.Header().Get("location"), "User should be redirected to url \"/register?IsRegisteringFailed=true\"")

	mockedUserContext.AssertExpectations(t)
}
