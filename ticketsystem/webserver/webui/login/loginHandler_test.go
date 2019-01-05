package login

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

/*
	A Valid Login Process should be possible.
*/
func TestUserLoginHandler_ServeHTTP_LoginSuccessful(t *testing.T) {
	userName := "TestUser"
	password := "TestPassword"
	token := "TestToken"
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := UserLoginHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
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
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, true, 1, "")
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

/*
	A failed login should return the same page with an error message.
*/
func TestUserLoginHandler_ServeHTTP_LoginFailed(t *testing.T) {
	userName := "TestUser"
	password := "TestPassword"
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	testee := UserLoginHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
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
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, true, 1, "")
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
