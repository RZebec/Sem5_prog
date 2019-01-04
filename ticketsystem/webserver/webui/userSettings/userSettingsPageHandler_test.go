package userSettings

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
func TestUserSettingsPageHandler_ServeHTTP__ValidRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/user_settings", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	rr := httptest.NewRecorder()

	testee := UserSettingsPageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	Only GET request should be possible.
*/
func TestUserSettingsPageHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	rr := httptest.NewRecorder()

	testee := UserSettingsPageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 405, rr.Code, "Status code 405 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	A error from rendering the template should result in a 500.
*/
func TestUserSettingsPageHandler_ServeHTTP_ContextError(t *testing.T) {
	req, err := http.NewRequest("GET", "/user_settings", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("TestError"))

	rr := httptest.NewRecorder()

	testee := UserSettingsPageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	Page should be returned with valid query parameter set.
*/
func TestUserSettingsPageHandler_ServeHTTP_ChangePasswordFailed(t *testing.T) {
	req, err := http.NewRequest("GET", "/user_settings?IsChangeFailed=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)

	testLogger := testhelpers.GetTestLogger()

	data := userSettingsPageData{
		IsChangeFailed: true,
	}

	data.UserIsAuthenticated = wrappers.IsAuthenticated(req.Context())
	data.UserIsAdmin = wrappers.IsAdmin(req.Context())
	data.Active = "settings"

	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "UserSettingsPage", data).Return(nil)

	testee := UserSettingsPageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}
