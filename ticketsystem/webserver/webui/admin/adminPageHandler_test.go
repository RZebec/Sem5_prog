package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	A logger for tests.
*/
func getTestLogger() logging.Logger {
	return logging.ConsoleLogger{SetTimeStamp: false}
}

func TestAdminPageHandler_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := getTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("GetIncomingMailApiKey").Return("1234")

	mockedApiContext.On("GetOutgoingMailApiKey").Return("4321")

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedUserContext.On("GetAllLockedUsers").Return([]user.User{{"Test@Test.de", 1,
	"Test","Test", user.RegisteredUser, user.WaitingToBeUnlocked}})

	rr := httptest.NewRecorder()

	testee := AdminPageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager, ApiContext: mockedApiContext}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")
}

func TestAdminPageHandlerOnError_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := getTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("GetIncomingMailApiKey").Return("1234")

	mockedApiContext.On("GetOutgoingMailApiKey").Return("4321")

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testError := templateManager.NewError("Error")

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(testError)

	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedUserContext.On("GetAllLockedUsers").Return([]user.User{{"Test@Test.de", 1,
		"Test","Test", user.RegisteredUser, user.WaitingToBeUnlocked}})

	testee := AdminPageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager, ApiContext: mockedApiContext}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code 500 should be returned")
}

func TestAdminPageHandlerWrongRequest_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("POST", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := getTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("GetIncomingMailApiKey").Return("1234")

	mockedApiContext.On("GetOutgoingMailApiKey").Return("4321")

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedUserContext.On("GetAllLockedUsers").Return([]user.User{{"Test@Test.de", 1,
		"Test","Test", user.RegisteredUser, user.WaitingToBeUnlocked}})

	testee := AdminPageHandler{UserContext: mockedUserContext, Logger: testLogger, TemplateManager: mockedTemplateManager, ApiContext: mockedApiContext}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Status code 405 should be returned")
}