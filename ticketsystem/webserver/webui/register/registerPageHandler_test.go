package register

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
	A valid request should be possible.
*/
func TestPageHandler_ServeHTTP_ValidRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	rr := httptest.NewRecorder()

	testee := PageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

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
	A error from rendering the template should result in a 500.
*/
func TestPageHandler_ServeHTTP_ContextError(t *testing.T) {
	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("TestError"))

	rr := httptest.NewRecorder()

	testee := PageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	Page should be returned with valid query parameter set.
*/
func TestPageHandler_ServeHTTP_RegisteringFailed(t *testing.T) {
	req, err := http.NewRequest("GET", "/register?IsRegisteringFailed=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	data := registerPageData{
		IsRegisteringFailed: true,
	}

	data.UserIsAuthenticated = wrappers.IsAuthenticated(req.Context())
	data.UserIsAdmin = wrappers.IsAuthenticated(req.Context())
	data.Active = "register"

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, data).Return(nil)

	rr := httptest.NewRecorder()

	testee := PageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	The User should be redirected to the index page if he is already logged in.
	The User should log out before registering a new user.
*/
func TestPageHandler_ServeHTTP_UserAlreadyLoggedIn_Redirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	rr := httptest.NewRecorder()

	testee := PageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusSeeOther, rr.Code, "Status code 303 should be returned")
	assert.Equal(t, "/", rr.Header().Get("location"), "User should be redirected to url \"/\"")

	mockedTemplateManager.AssertExpectations(t)
}
