// 5894619, 6720876, 9793350
package userSettings

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
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

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")

	testLogger := testhelpers.GetTestLogger()

	testUser := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserById", 5).Return(true, testUser)

	data := userSettingsPageData{
		IsChangeFailed:   "NotSet",
		UserIsOnVacation: false,
	}

	data.UserIsAuthenticated = true
	data.UserIsAdmin = false
	data.Active = "settings"

	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, data).Return(nil)

	rr := httptest.NewRecorder()

	testee := SettingsPageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Only GET request should be possible.
*/
func TestUserSettingsPageHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/user_settings", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")

	testLogger := testhelpers.GetTestLogger()

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	rr := httptest.NewRecorder()

	testee := SettingsPageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 405, rr.Code, "Status code 405 should be returned")

	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A error from rendering the template should result in a 500.
*/
func TestUserSettingsPageHandler_ServeHTTP_ContextError(t *testing.T) {
	req, err := http.NewRequest("GET", "/user_settings", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")

	testLogger := testhelpers.GetTestLogger()

	testUser := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserById", 5).Return(true, testUser)

	data := userSettingsPageData{
		IsChangeFailed:   "NotSet",
		UserIsOnVacation: false,
	}

	data.UserIsAuthenticated = true
	data.UserIsAdmin = false
	data.Active = "settings"

	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, data).Return(errors.New("TestError"))

	rr := httptest.NewRecorder()

	testee := SettingsPageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")

	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Page should be returned with valid query parameter set.
*/
func TestUserSettingsPageHandler_ServeHTTP_ChangePasswordFailed(t *testing.T) {
	req, err := http.NewRequest("GET", "/user_settings?IsChangeFailed=yes", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")

	testLogger := testhelpers.GetTestLogger()

	testUser := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}

	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserById", 5).Return(true, testUser)

	data := userSettingsPageData{
		IsChangeFailed:   "yes",
		UserIsOnVacation: false,
	}

	data.UserIsAuthenticated = true
	data.UserIsAdmin = false
	data.Active = "settings"

	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, data).Return(nil)

	testee := SettingsPageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}
