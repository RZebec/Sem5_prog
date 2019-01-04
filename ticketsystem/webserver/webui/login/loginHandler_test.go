package login

import (
	"de/vorlesung/projekt/IIIDDD/shared"
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
	A GET request for the login page should be possible
*/
func TestLoginHandler_ServeHTTPGetLoginPage_UserNotLoggedIn(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := LoginHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
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
	handler := http.HandlerFunc(testee.ServeHTTPGetLoginPage)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return status code 200")
	assert.Equal(t, "", newLocation, "Should not be redirected")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

/*
	Only GET request should be possible.
*/
func TestLoginHandler_ServeHTTPGetLoginPage_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	rr := httptest.NewRecorder()

	testee := LoginHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTPGetLoginPage)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 405, rr.Code, "Status code 405 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	An error from the template manager should return a 500.
*/
func TestLoginHandler_ServeHTTPGetLoginPage_ContextError(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := LoginHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
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
	handler := http.HandlerFunc(testee.ServeHTTPGetLoginPage)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/", newLocation, "Should redirect to \"/\"")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

/*
	A already logged in user should be redirected to the index page.
*/
func TestLoginHandler_ServeHTTPGetLoginPage_UserAlreadyLoggedIn(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := LoginHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTPGetLoginPage)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode, "Should return status code 303")
	assert.Equal(t, "/", newLocation, "Should be redirected to /")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

func TestLoginHandler_ServeHTTPPostLoginData_LoginSuccessful(t *testing.T) {
	userName := "TestUser"
	password := "TestPassword"
	token := "TestToken"
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := LoginHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	mockedUserContext.On("Login", userName, password).Return(true, token, nil)

	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Form = url.Values{}
	req.Form.Add("userName", userName)
	req.Form.Add("password", password)

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTPPostLoginData)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, true, 1,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, 302, resp.StatusCode, "Should return status code 302")
	assert.Equal(t, "/", newLocation, "Should be redirected to /")
	cookieExists, cookieValue := testhelpers.GetCookieValue(resp.Cookies(), shared.AccessTokenCookieName)
	assert.True(t, cookieExists, "The cookie should be set")
	assert.Equal(t, token, cookieValue, "The cookie should be set to the correct value")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

func TestLoginHandler_ServeHTTPPostLoginData_LoginFailed(t *testing.T) {
	userName := "TestUser"
	password := "TestPassword"
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := LoginHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	mockedUserContext.On("Login", userName, password).Return(false, "", nil)

	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Form = url.Values{}
	req.Form.Add("userName", userName)
	req.Form.Add("password", password)

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTPPostLoginData)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, true, 1,"")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, 302, resp.StatusCode, "Should return status code 302")
	assert.Equal(t, "/login?IsLoginFailed=true", newLocation, "Should be redirected to /login?IsLoginFailed=true")
	cookieExists, _ := testhelpers.GetCookieValue(resp.Cookies(), shared.AccessTokenCookieName)
	assert.False(t, cookieExists, "The cookie should not be set")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}
