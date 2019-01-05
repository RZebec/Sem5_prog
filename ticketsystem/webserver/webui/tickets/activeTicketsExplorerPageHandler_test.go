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
	Get test tickets with different states.
*/
func getTestTickets() []ticketData.TicketInfo {
	testEditor := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}
	testCreator := ticketData.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	var ticketInfos []ticketData.TicketInfo
	ticketInfos = append(ticketInfos, ticketData.TicketInfo{Id: 1, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator,
		CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticketData.Open})
	ticketInfos = append(ticketInfos, ticketData.TicketInfo{Id: 1, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator,
		CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticketData.Processing})
	ticketInfos = append(ticketInfos, ticketData.TicketInfo{Id: 1, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator,
		CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticketData.Processing})
	ticketInfos = append(ticketInfos, ticketData.TicketInfo{Id: 1, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator,
		CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticketData.Closed})

	return ticketInfos
}

/*
	A Valid Http Request.
*/
func TestActiveTicketsExplorerPageHandler_ServeHTTP_ValidRequest(t *testing.T) {
	userIsAuthenticated := false
	userIsAdmin := false
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := ActiveTicketsExplorerPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	// Only tickets with state processing should be shown.
	testTickets := getTestTickets()
	var expectedTicketData []ticketData.TicketInfo
	for _, possibleTicket := range testTickets {
		if possibleTicket.State == ticketData.Processing {
			expectedTicketData = append(expectedTicketData, possibleTicket)
		}
	}
	expectedPageData := activeTicketsExplorerPageData{
		Tickets: expectedTicketData,
	}
	expectedPageData.UserIsAdmin = userIsAdmin
	expectedPageData.UserIsAuthenticated = userIsAuthenticated
	expectedPageData.Active = "active_tickets"

	mockedTicketContext.On("GetAllTicketInfo").Return(testTickets)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketExplorerPage", expectedPageData).Return(nil)

	req, err := http.NewRequest("GET", "/active_tickets", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), userIsAuthenticated, userIsAdmin, -1, "")
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
func TestActiveTicketsExplorerPageHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := ActiveTicketsExplorerPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("POST", "/active_tickets", nil)
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
func TestActiveTicketsExplorerPageHandler_ServeHTTP_ContextError_RenderError(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := ActiveTicketsExplorerPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	testEditor := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}
	testCreator := ticketData.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	testTickets := []ticketData.TicketInfo{{1, "TicketTest", testEditor, true, testCreator, time.Now(), time.Now(), ticketData.Open}}

	mockedTicketContext.On("GetAllTicketInfo").Return(testTickets)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketExplorerPage", mock.Anything).Return(errors.New("TestError"))

	req, err := http.NewRequest("GET", "/active_tickets", nil)
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
