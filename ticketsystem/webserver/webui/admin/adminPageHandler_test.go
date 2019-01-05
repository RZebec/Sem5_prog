package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	A valid request should be possible.
*/
func TestAdminPageHandler_ServeHTTP_ValidRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("GetIncomingMailApiKey").Return("1234")

	mockedApiContext.On("GetOutgoingMailApiKey").Return("4321")

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedUserContext.On("GetAllLockedUsers").Return([]userData.User{{"Test@Test.de", 1,
		"Test", "Test", userData.RegisteredUser, userData.WaitingToBeUnlocked}})

	rr := httptest.NewRecorder()

	testee := PageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager, ApiContext: mockedApiContext}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedApiContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A error from rendering the template should result in a 500.
*/
func TestAdminPageHandler_ServeHTTP_RenderTemplateError_500Returned(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("GetIncomingMailApiKey").Return("1234")

	mockedApiContext.On("GetOutgoingMailApiKey").Return("4321")

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testError := templateManager.NewError("Error")

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(testError)

	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedUserContext.On("GetAllLockedUsers").Return([]userData.User{{"Test@Test.de", 1,
		"Test", "Test", userData.RegisteredUser, userData.WaitingToBeUnlocked}})

	testee := PageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager, ApiContext: mockedApiContext}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code 500 should be returned")
}

/*
	Only GET Methods should be allowed.
*/
func TestAdminPageHandler_ServeHTTP_WrongRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("GetIncomingMailApiKey").Return("1234")

	mockedApiContext.On("GetOutgoingMailApiKey").Return("4321")

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedUserContext.On("GetAllLockedUsers").Return([]userData.User{{"Test@Test.de", 1,
		"Test", "Test", userData.RegisteredUser, userData.WaitingToBeUnlocked}})

	testee := PageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager, ApiContext: mockedApiContext}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Status code 405 should be returned")
}

/*
	Failed change should be shown.
*/
func TestAdminPageHandler_ServeHTTP_FailedChange(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin?IsChangeFailed=yes", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	lockedUsers := []userData.User{{"Test@Test.de", 1,
		"Test", "Test", userData.RegisteredUser, userData.WaitingToBeUnlocked}}

	data := adminPageData{
		Users:              lockedUsers,
		IncomingMailApiKey: "1234",
		OutgoingMailApiKey: "4321",
		IsChangeFailed:     "yes",
	}
	data.UserIsAdmin = true
	data.UserIsAuthenticated = true
	data.Active = "admin"

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("GetIncomingMailApiKey").Return("1234")

	mockedApiContext.On("GetOutgoingMailApiKey").Return("4321")

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, data).Return(nil)

	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedUserContext.On("GetAllLockedUsers").Return(lockedUsers)

	rr := httptest.NewRecorder()

	testee := PageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager, ApiContext: mockedApiContext}

	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, true, 1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedApiContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}
