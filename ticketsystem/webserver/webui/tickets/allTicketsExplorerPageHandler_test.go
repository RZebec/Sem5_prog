// 5894619, 6720876, 9793350
package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
	A Valid Http Request.
*/
func TestAllTicketsExplorerPageHandler_ServeHTTP_ValidRequest(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := AllTicketsExplorerPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	testEditor := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}
	testCreator := ticketData.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	testTickets := []ticketData.TicketInfo{{1, "TicketTest", testEditor, true, testCreator, time.Now(), time.Now(), ticketData.Open}}

	mockedTicketContext.On("GetAllTicketInfo").Return(testTickets)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketExplorerPage", mock.Anything).Return(nil)

	req, err := http.NewRequest("GET", "/all_tickets", nil)
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

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

/*
	Only GET request should be possible.
*/
func TestAllTicketsExplorerPageHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := AllTicketsExplorerPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("POST", "/all_tickets", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Should return status code 405")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

/*
	A error from rendering the template should result in a 500.
*/
func TestAllTicketsExplorerPageHandler_ServeHTTP_ContextError_RenderError(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := AllTicketsExplorerPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	testEditor := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}
	testCreator := ticketData.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	testTickets := []ticketData.TicketInfo{{1, "TicketTest", testEditor, true, testCreator, time.Now(), time.Now(), ticketData.Open}}

	mockedTicketContext.On("GetAllTicketInfo").Return(testTickets)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketExplorerPage", mock.Anything).Return(errors.New("TestError"))

	req, err := http.NewRequest("GET", "/all_tickets", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return status code 500")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}
