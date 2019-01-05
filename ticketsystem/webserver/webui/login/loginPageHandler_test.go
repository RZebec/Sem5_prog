package login

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	A GET request for the login page should be possible
*/
func TestPageHandler_ServeHTTP_UserNotLoggedIn(t *testing.T) {
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := PageHandler{TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	expectedLoginData := loginPageData{
		IsLoginFailed: false,
	}
	expectedLoginData.UserIsAuthenticated = false
	expectedLoginData.UserIsAdmin = false
	expectedLoginData.Active = "login"
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "LoginPage", expectedLoginData).Return(nil)

	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return status code 200")
	assert.Equal(t, "", newLocation, "Should not be redirected")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	Only GET request should be possible.
*/
func TestPageHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	rr := httptest.NewRecorder()

	testee := PageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 405, rr.Code, "Status code 405 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	An error from the template manager should return a 500.
*/
func TestPageHandler_ServeHTTP_ContextError(t *testing.T) {
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := PageHandler{TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	expectedLoginData := loginPageData{
		IsLoginFailed: false,
	}
	expectedLoginData.UserIsAuthenticated = false
	expectedLoginData.UserIsAdmin = false
	expectedLoginData.Active = "login"
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "LoginPage", expectedLoginData).Return(errors.New("TestError"))

	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/", newLocation, "Should redirect to \"/\"")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	A already logged in user should be redirected to the index page.
*/
func TestPageHandler_ServeHTTP_UserAlreadyLoggedIn(t *testing.T) {
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := PageHandler{TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode, "Should return status code 303")
	assert.Equal(t, "/", newLocation, "Should be redirected to /")

	mockedTemplateManager.AssertExpectations(t)
}
