package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"github.com/stretchr/testify/assert"
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

	templateManager.LoadTemplates(testLogger)

	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedUserContext.On("GetAllLockedUsers").Return([]user.User{{"Test@Test.de", 1,
	"Test","Test", user.RegisteredUser, user.WaitingToBeUnlocked}})

	rr := httptest.NewRecorder()

	testee := AdminPageHandler{UserContext: mockedUserContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")
}