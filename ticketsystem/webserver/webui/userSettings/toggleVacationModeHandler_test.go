package userSettings

import (
	"de/vorlesung/projekt/IIIDDD/shared"
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
func TestToggleVacationModeHandler_ServeHTTP_ValidRequest_RedirectedToUserSettings(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("vacationMode", "true")

	cookie := http.Cookie{Name: shared.AccessTokenCookieName, Value: "test"}
	req.AddCookie(&cookie)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("EnableVacationMode", "test").Return(nil)

	testee := ToggleVacationModeHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 302, rr.Code, "Status code 302 should be returned")
	assert.Equal(t, "/user_settings", rr.Header().Get("location"), "User should be redirected to url \"/user_settings\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	A valid request should be possible.
*/
func TestToggleVacationModeHandler_ServeHTTP_ValidRequest_DisableVacationMode_RedirectedToUserSettings(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("vacationMode", "false")

	cookie := http.Cookie{Name: shared.AccessTokenCookieName, Value: "test"}
	req.AddCookie(&cookie)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("DisableVacationMode", "test").Return(nil)

	testee := ToggleVacationModeHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 302, rr.Code, "Status code 302 should be returned")
	assert.Equal(t, "/user_settings", rr.Header().Get("location"), "User should be redirected to url \"/user_settings\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	Only POST Method should be able.
*/
func TestToggleVacationModeHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("vacationMode", "true")

	cookie := http.Cookie{Name: shared.AccessTokenCookieName, Value: "test"}
	req.AddCookie(&cookie)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	testee := ToggleVacationModeHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 405, rr.Code, "Status code 405 should be returned")

	mockedUserContext.AssertExpectations(t)
}

/*
	Parsing Error should return a 400.
*/
func TestToggleVacationModeHandler_ServeHTTP_ParsingError(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("vacationMode", "fasfasf")

	cookie := http.Cookie{Name: shared.AccessTokenCookieName, Value: "test"}
	req.AddCookie(&cookie)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := ToggleVacationModeHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code 400 should be returned")
	assert.Equal(t, "/", rr.Header().Get("location"), "User should be redirected to url \"/\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	Cookie Error should return a 400.
*/
func TestToggleVacationModeHandler_ServeHTTP_CookieError(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("vacationMode", "true")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := ToggleVacationModeHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code 400 should be returned")
	assert.Equal(t, "/", rr.Header().Get("location"), "User should be redirected to url \"/\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	Enable Vacation Mode error should result in a 400.
*/
func TestToggleVacationModeHandler_ServeHTTP_EnableVacationModeError(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("vacationMode", "true")

	cookie := http.Cookie{Name: shared.AccessTokenCookieName, Value: "test"}
	req.AddCookie(&cookie)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("EnableVacationMode", "test").Return(errors.New("TestError"))

	testee := ToggleVacationModeHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code 400 should be returned")
	assert.Equal(t, "/", rr.Header().Get("location"), "User should be redirected to url \"/\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	Disable Vacation Mode error should result in a 400.
*/
func TestToggleVacationModeHandler_ServeHTTP_DisableVacationModeError(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	req.Form = url.Values{}
	req.Form.Add("vacationMode", "false")

	cookie := http.Cookie{Name: shared.AccessTokenCookieName, Value: "test"}
	req.AddCookie(&cookie)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("DisableVacationMode", "test").Return(errors.New("TestError"))

	testee := ToggleVacationModeHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code 400 should be returned")
	assert.Equal(t, "/", rr.Header().Get("location"), "User should be redirected to url \"/\"")

	mockedUserContext.AssertExpectations(t)
}