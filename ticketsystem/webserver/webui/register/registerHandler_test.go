package register

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
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
func TestRegisterHandler_ServeHTTPPostRegisteringData_ValidRequest_RedirectedToLoginPage(t *testing.T) {
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

	testee := RegisterHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTPPostRegisteringData)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 303, rr.Code, "Status code 303 should be returned")
	assert.Equal(t, "/login", rr.Header().Get("location"), "User should be redirected to url \"/login\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	Only POST request should be possible.
*/
func TestRegisterHandler_ServeHTTPPostRegisteringData_WrongRequestMethod(t *testing.T) {
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

	testee := RegisterHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTPPostRegisteringData)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Status code 405 should be returned")

	mockedUserContext.AssertExpectations(t)
}

/*
	An unsuccessful registering procedure should return the same page.
*/
func TestRegisterHandler_ServeHTTPPostRegisteringData_UnsuccessfulRegistering_RedirectedToSamePageWithMessage(t *testing.T) {
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

	testee := RegisterHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTPPostRegisteringData)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	rr.Result()

	assert.Equal(t, 303, rr.Code, "Status code 303 should be returned")
	assert.Equal(t, "/register?IsRegisteringFailed=true", rr.Header().Get("location"), "User should be redirected to url \"/register?IsRegisteringFailed=true\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	An unsuccessful registering procedure should return the same page.
*/
func TestRegisterHandler_ServeHTTPPostRegisteringData_ContextReturnError_RedirectedToSamePageWithMessage(t *testing.T) {
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

	testee := RegisterHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTPPostRegisteringData)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	rr.Result()

	assert.Equal(t, 303, rr.Code, "Status code 303 should be returned")
	assert.Equal(t, "/register?IsRegisteringFailed=true", rr.Header().Get("location"), "User should be redirected to url \"/register?IsRegisteringFailed=true\"")

	mockedUserContext.AssertExpectations(t)
}

/*
	A valid request should be possible.
*/
func TestRegisterHandler_ServeHTTPGetRegisterPage_ValidRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	rr := httptest.NewRecorder()

	testee := RegisterHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTPGetRegisterPage)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	Only GET request should be possible.
*/
func TestRegisterHandler_ServeHTTPGetRegisterPage_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	rr := httptest.NewRecorder()

	testee := RegisterHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTPGetRegisterPage)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 405, rr.Code, "Status code 405 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	A error from rendering the template should result in a 500.
*/
func TestRegisterHandler_ServeHTTPGetRegisterPage_ContextError(t *testing.T) {
	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("TestError"))

	rr := httptest.NewRecorder()

	testee := RegisterHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTPGetRegisterPage)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	Page should be returned with valid query parameter set.
*/
func TestRegisterHandler_ServeHTTPGetRegisterPage_RegisteringFailed(t *testing.T) {
	req, err := http.NewRequest("GET", "/register?IsRegisteringFailed=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	rr := httptest.NewRecorder()

	testee := RegisterHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTPGetRegisterPage)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	The User should be redirected to the index page if he is already logged in.
	The User should log out before registering a new user.
*/
func TestRegisterHandler_ServeHTTPGetRegisterPage_UserAlreadyLoggedIn_Redirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	rr := httptest.NewRecorder()

	testee := RegisterHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTPGetRegisterPage)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusSeeOther, rr.Code, "Status code 303 should be returned")
	assert.Equal(t, "/", rr.Header().Get("location"), "User should be redirected to url \"/\"")

	mockedTemplateManager.AssertExpectations(t)
}
