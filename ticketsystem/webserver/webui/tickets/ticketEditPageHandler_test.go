package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
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
func TestTicketEditPageHandler_ServeHTTP_ValidRequest(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := TicketEditPageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	testEditor := user.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: user.RegisteredUser, State: user.Active}

	testCreator := ticket.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	testTicketInfo := ticket.TicketInfo{Id: 5, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State:  ticket.Processing}
	testMessages := []ticket.MessageEntry{{Id: 0, CreatorMail: "test@test.de", Content: "TestContent2", OnlyInternal: false}}

	testTicket := ticket.CreateTestTicket(testTicketInfo, testMessages)

	testUser := user.User{Mail: "Test25@Test.de", UserId: 25, FirstName: "Dieter22", LastName: "Dietrich", Role: user.RegisteredUser, State: user.Active}

	testCreator2 := ticket.Creator{Mail: "Ivan@Test.de", FirstName: "Ivan", LastName: "Muller"}
	testTicketInfo2 := ticket.TicketInfo{Id: 0, Title: "TicketTest2", Editor: user.GetInvalidDefaultUser(), HasEditor: false, Creator: testCreator2, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Open}

	testCreator3 := ticket.Creator{Mail: "Ivan2@Test.de", FirstName: "Ivan2", LastName: "Muller"}
	testTicketInfo3 := ticket.TicketInfo{Id: 1, Title: "TicketTest3", Editor: user.GetInvalidDefaultUser(), HasEditor: false, Creator: testCreator3, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Closed}

	testEditor4 := user.User{Mail: "Test44@Test.de", UserId: 44, FirstName: "Dieter33", LastName: "Dietrich", Role: user.RegisteredUser, State: user.Active}
	testCreator4 := ticket.Creator{Mail: "Ivan3@Test.de", FirstName: "Ivan3", LastName: "Muller"}
	testTicketInfo4 := ticket.TicketInfo{Id: 2, Title: "TicketTest4", Editor: testEditor4, HasEditor: true, Creator: testCreator4, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Processing}

	req, err := http.NewRequest("GET", "/ticket_edit/5", nil)
	if err != nil {
		t.Fatal(err)
	}

	allTestTickets := []ticket.TicketInfo{testTicketInfo, testTicketInfo2, testTicketInfo3, testTicketInfo4}
	allTestUsers := []user.User{testEditor, testUser, testEditor4}
	filteredTickets := []ticket.TicketInfo{testTicketInfo2, testTicketInfo3, testTicketInfo4}

	data := ticketEditPageData{
		TicketInfo: 	testTicketInfo,
		OtherTickets:	filteredTickets,
		Users:			allTestUsers,
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	data.UserIsAdmin = false
	data.UserIsAuthenticated = true
	data.Active = "ticket_edit"

	mockedTicketContext.On("GetTicketById", 5).Return(true, testTicket)
	mockedTicketContext.On("GetAllTicketInfo").Return(allTestTickets)

	mockedUserContext.On( "GetAllActiveUsers").Return(allTestUsers)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketEditPage", data).Return(nil)


	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return status code 200")
	assert.Equal(t, "", newLocation, "Should not be redirected")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Only GET requests should be possible.
*/
func TestTicketEditPageHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := TicketEditPageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("POST", "/ticket_edit/5", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Should return status code 405")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	ID Conversion Errors should return a 400.
*/
func TestTicketEditPageHandler_ServeHTTP_ConversionError(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := TicketEditPageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("GET", "/ticket_edit/aaaa", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return status code 400")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	If the ticket doesnt exist it should return a 404.
*/
func TestTicketEditPageHandler_ServeHTTP_TicketDoesNotExist(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := TicketEditPageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("GET", "/ticket_edit/5", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockedTicketContext.On("GetTicketById", 5).Return(false, new(ticket.Ticket))

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Should return status code 404")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Template render error should return a 500.
*/
func TestTicketEditPageHandler_ServeHTTP_RenderTemplateError(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := TicketEditPageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	testEditor := user.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: user.RegisteredUser, State: user.Active}

	testCreator := ticket.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	testTicketInfo := ticket.TicketInfo{Id: 5, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State:  ticket.Processing}
	testMessages := []ticket.MessageEntry{{Id: 0, CreatorMail: "test@test.de", Content: "TestContent2", OnlyInternal: false}}

	testTicket := ticket.CreateTestTicket(testTicketInfo, testMessages)

	testUser := user.User{Mail: "Test25@Test.de", UserId: 25, FirstName: "Dieter22", LastName: "Dietrich", Role: user.RegisteredUser, State: user.Active}

	testCreator2 := ticket.Creator{Mail: "Ivan@Test.de", FirstName: "Ivan", LastName: "Muller"}
	testTicketInfo2 := ticket.TicketInfo{Id: 0, Title: "TicketTest2", Editor: user.GetInvalidDefaultUser(), HasEditor: false, Creator: testCreator2, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Open}

	testCreator3 := ticket.Creator{Mail: "Ivan2@Test.de", FirstName: "Ivan2", LastName: "Muller"}
	testTicketInfo3 := ticket.TicketInfo{Id: 1, Title: "TicketTest3", Editor: user.GetInvalidDefaultUser(), HasEditor: false, Creator: testCreator3, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Closed}

	testEditor4 := user.User{Mail: "Test44@Test.de", UserId: 44, FirstName: "Dieter33", LastName: "Dietrich", Role: user.RegisteredUser, State: user.Active}
	testCreator4 := ticket.Creator{Mail: "Ivan3@Test.de", FirstName: "Ivan3", LastName: "Muller"}
	testTicketInfo4 := ticket.TicketInfo{Id: 2, Title: "TicketTest4", Editor: testEditor4, HasEditor: true, Creator: testCreator4, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Processing}

	req, err := http.NewRequest("GET", "/ticket_edit/5", nil)
	if err != nil {
		t.Fatal(err)
	}

	allTestTickets := []ticket.TicketInfo{testTicketInfo, testTicketInfo2, testTicketInfo3, testTicketInfo4}
	allTestUsers := []user.User{testEditor, testUser, testEditor4}
	filteredTickets := []ticket.TicketInfo{testTicketInfo2, testTicketInfo3, testTicketInfo4}

	data := ticketEditPageData{
		TicketInfo: 	testTicketInfo,
		OtherTickets:	filteredTickets,
		Users:			allTestUsers,
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	data.UserIsAdmin = false
	data.UserIsAuthenticated = true
	data.Active = "ticket_edit"

	mockedTicketContext.On("GetTicketById", 5).Return(true, testTicket)
	mockedTicketContext.On("GetAllTicketInfo").Return(allTestTickets)

	mockedUserContext.On( "GetAllActiveUsers").Return(allTestUsers)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketEditPage", data).Return(errors.New("TestError"))


	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/", newLocation, "Should be redirected to \"/\"")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}
