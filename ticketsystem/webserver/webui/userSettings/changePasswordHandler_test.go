package userSettings

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

/*
	A valid request should be possible.
*/
func TestChangePasswordHandler_ServeHTTP_ValidRequest_RedirectedToUserSettings(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("old_password", "aaBB11==")
	req.Form.Add("new_password", "bbCC22==")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("ChangePassword", "", "aaBB11==", "bbCC22==").Return(true, nil)

	testee := ChangePasswordHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 302, rr.Code, "Status code 302 should be returned")
	assert.Equal(t, "/user_settings", rr.Header().Get("location"), "User should be redirected to url \"/user_settings\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	Only POST request should be possible.
*/
func TestChangePasswordHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("old_password", "aaBB11==")
	req.Form.Add("new_password", "bbCC22==")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := ChangePasswordHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Status code 405 should be returned")

	mockedUserContext.AssertExpectations(t)
}

/*
	An unsuccessful change procedure should return the same page.
*/
func TestChangePasswordHandler_ServeHTTP_UnsuccessfulChange_RedirectedToSamePageWithMessage(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("old_password", "aaBB11==")
	req.Form.Add("new_password", "bbCC22==")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("ChangePassword","", "aaBB11==", "bbCC22==").Return(false, nil)

	testee := ChangePasswordHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 302, rr.Code, "Status code 302 should be returned")
	assert.Equal(t, "/user_settings?IsChangeFailed=true", rr.Header().Get("location"), "User should be redirected to url \"/user_register?IsChangeFailed=true\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	An unsuccessful registering procedure should return the same page.
*/
func TestChangePasswordHandler_ServeHTTP_ContextReturnError_RedirectedToSamePageWithMessage(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_register", nil)
	req.Form = url.Values{}
	req.Form.Add("old_password", "aaBB11==")
	req.Form.Add("new_password", "bbCC22==")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("ChangePassword", "", "aaBB11==", "bbCC22==").Return(false, errors.New("TestError"))

	testee := ChangePasswordHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 302, rr.Code, "Status code 302 should be returned")
	assert.Equal(t, "/user_settings?IsChangeFailed=true", rr.Header().Get("location"), "User should be redirected to url \"/user_register?IsChangeFailed=true\"")

	mockedUserContext.AssertExpectations(t)
}
