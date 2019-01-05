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
func TestTicketViewPageHandler_ServeHTTP_ValidRequest(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketViewPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager, UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger()}

	testEditor := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}
	testCreator := ticketData.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	testTicketInfo := ticketData.TicketInfo{Id: 5, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticketData.Open}
	testMessages := []ticketData.MessageEntry{{Id: 0, CreatorMail: "test@test.de", Content: "TestContent2", OnlyInternal: false}}

	testTicket := ticketData.CreateTestTicket(testTicketInfo, testMessages)

	mockedTicketContext.On("GetTicketById", 5).Return(true, testTicket)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketViewPage", mock.Anything).Return(nil)

	req, err := http.NewRequest("GET", "/ticketData/5", nil)
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
	mockedUserContext.AssertExpectations(t)
}

/*
	Only GET request should be possible.
*/
func TestTicketViewPageHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketViewPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager, UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("POST", "/ticket/5", nil)
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
	mockedUserContext.AssertExpectations(t)
}

/*
	A error from rendering the template should result in a 500.
*/
func TestTicketViewPageHandler_ServeHTTP_ContextError_RenderError(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketViewPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager, UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger()}

	testEditor := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}
	testCreator := ticketData.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	testTicketInfo := ticketData.TicketInfo{Id: 5, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticketData.Open}
	testMessages := []ticketData.MessageEntry{{Id: 0, CreatorMail: "test@test.de", Content: "TestContent2", OnlyInternal: false}}

	testTicket := ticketData.CreateTestTicket(testTicketInfo, testMessages)

	mockedTicketContext.On("GetTicketById", 5).Return(true, testTicket)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketViewPage", mock.Anything).Return(errors.New("TestError"))

	req, err := http.NewRequest("GET", "/ticket/5", nil)
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
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/", newLocation, "Should be redirected to \"/\"")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A error from id conversion should result in a 400.
*/
func TestTicketViewPageHandler_ServeHTTP_IdConversionError(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketViewPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager, UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("GET", "/ticket/asdh", nil)
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
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return status code 400")
	assert.Equal(t, "/", newLocation, "Should be redirected to \"/\"")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Ticket doesn´t exist should result in a 404.
*/
func TestTicketViewPageHandler_ServeHTTP_TicketDoesNotExist(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketViewPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager, UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger()}

	mockedTicketContext.On("GetTicketById", 5).Return(false, new(ticketData.Ticket))

	req, err := http.NewRequest("GET", "/ticket/5", nil)
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
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Should return status code 404")
	assert.Equal(t, "/", newLocation, "Should be redirected to \"/\"")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A User that is not logged in should not see internal messages.
*/
func TestTicketViewPageHandler_ServeHTTP_DoNotShowInternalOnlyMessages(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketViewPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager, UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger()}

	testEditor := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}
	testCreator := ticketData.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	testTicketInfo := ticketData.TicketInfo{Id: 5, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticketData.Open}
	testMessages := []ticketData.MessageEntry{{Id: 1, CreatorMail: "test1@test.de", Content: "TestContent1", OnlyInternal: false}, {Id: 2, CreatorMail: "test2@test.de", Content: "TestContent2", OnlyInternal: true}}

	testTicket := ticketData.CreateTestTicket(testTicketInfo, testMessages)

	externalTestMessages := []ticketData.MessageEntry{{Id: 1, CreatorMail: "test1@test.de", Content: "TestContent1", OnlyInternal: false}}

	req, err := http.NewRequest("GET", "/ticket/5", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")

	data := ticketViewPageData{
		TicketInfo: testTicketInfo,
		Messages:   externalTestMessages,
	}

	data.UserIsAdmin = false
	data.UserIsAuthenticated = false
	data.Active = "all_tickets"

	mockedTicketContext.On("GetTicketById", 5).Return(true, testTicket)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketViewPage", data).Return(nil)

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
	A User that is logged in should see all messages.
*/
func TestTicketViewPageHandler_ServeHTTP_ShowInternalOnlyMessages(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketViewPageHandler{TicketContext: mockedTicketContext, TemplateManager: mockedTemplateManager, UserContext: mockedUserContext,
		Logger: testhelpers.GetTestLogger()}

	testEditor := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}
	testCreator := ticketData.Creator{Mail: "Test@Test.de", FirstName: "Max", LastName: "Muller"}
	testTicketInfo := ticketData.TicketInfo{Id: 5, Title: "TicketTest", Editor: testEditor, HasEditor: true, Creator: testCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticketData.Open}
	testMessages := []ticketData.MessageEntry{{Id: 1, CreatorMail: "test1@test.de", Content: "TestContent1", OnlyInternal: false}, {Id: 2, CreatorMail: "test2@test.de", Content: "TestContent2", OnlyInternal: true}}

	testTicket := ticketData.CreateTestTicket(testTicketInfo, testMessages)

	req, err := http.NewRequest("GET", "/ticket/5", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")

	data := ticketViewPageData{
		TicketInfo: testTicketInfo,
		Messages:   testMessages,
		UserName:   testEditor.Mail,
	}

	data.UserIsAdmin = false
	data.UserIsAuthenticated = true
	data.Active = "all_tickets"

	mockedTicketContext.On("GetTicketById", 5).Return(true, testTicket)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketViewPage", data).Return(nil)
	mockedUserContext.On("GetUserById", 5).Return(true, testEditor)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return status code 200")
	assert.Equal(t, "", newLocation, "Should not be redirected")

	mockedTicketContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}
