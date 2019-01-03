package login

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLoginHandler_ServeHTTPGetLoginPage(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := LoginHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	expectedLoginData := loginPageData{
		IsLoginFailed: false,
	}
	expectedLoginData.UserIsAuthenticated = true
	expectedLoginData.UserIsAdmin = true
	expectedLoginData.Active = "login"
	mockedTemplateManager.On("RenderTemplate",mock.Anything, "LoginPage", expectedLoginData ).Return(nil)

	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTPGetLoginPage)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, true)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode, "Should return status code 303")
	assert.Equal(t, "/", newLocation , "Should be redirected to /")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

func TestLoginHandler_ServeHTTPGetLoginPage_LoginRedirected(t *testing.T) {
	userName := "TestUser"
	password := "TestPassword"
	token := "TestToken"
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := LoginHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	expectedLoginData := loginPageData{
		IsLoginFailed: false,
	}
	expectedLoginData.UserIsAuthenticated = true
	expectedLoginData.UserIsAdmin = true
	expectedLoginData.Active = "login"
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
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, true)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, 302, resp.StatusCode, "Should return status code 302")
	assert.Equal(t, "/", newLocation , "Should be redirected to /")
	cookieExists, cookieValue := testhelpers.GetCookieValue(resp.Cookies(), shared.AccessTokenCookieName)
	assert.True(t, cookieExists, "The cookie should be set")
	assert.Equal(t, token, cookieValue, "The cookie should be set to the correct value")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

