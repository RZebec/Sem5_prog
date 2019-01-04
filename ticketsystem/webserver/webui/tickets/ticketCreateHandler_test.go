package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

/*
	A Valid Http Request for logged In user.
*/
func TestTicketCreateHandler_ServeHTTP_ValidRequestLoggedIn(t *testing.T) {
	mail := "test@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := "false"

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)

	testUser := user.User{Mail: mail, UserId: 5, FirstName: "Max", LastName: "Muller", Role: user.RegisteredUser, State: user.Active}
	testMessage := ticket.MessageEntry{Id: 0, CreatorMail: mail, Content: message, OnlyInternal: false}

	testTicketCreator := ticket.Creator{Mail: mail, FirstName: "Max", LastName: "Muller"}

	testTicketInfo := ticket.TicketInfo{Id: 26, Title: title, Editor: testUser, HasEditor: true, Creator: testTicketCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Processing}

	ticket := ticket.CreateTestTicket(testTicketInfo, []ticket.MessageEntry{testMessage})

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicketForInternalUser", title, testUser, testMessage).Return(ticket, nil)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", mail).Return(true, testUser.UserId)
	mockedUserContext.On("GetUserById", testUser.UserId).Return(true, testUser)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, 302, resp.StatusCode, "Should return status code 302")
	assert.Equal(t, "/ticket/26", newLocation, "Should be redirected to \"/ticket/26\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A Valid Http Request for not logged In user.
*/
func TestTicketCreateHandler_ServeHTTP_ValidRequestNotLoggedIn(t *testing.T) {
	mail := "test@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := "false"
	firstName := "Max"
	lastName := "Muller"

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)
	req.Form.Add("first_name", firstName)
	req.Form.Add("last_name", lastName)

	testMessage := ticket.MessageEntry{Id: 0, CreatorMail: mail, Content: message, OnlyInternal: false}

	testTicketCreator := ticket.Creator{Mail: mail, FirstName: firstName, LastName: lastName}

	testTicketInfo := ticket.TicketInfo{Id: 26, Title: title, Editor: user.GetInvalidDefaultUser(), HasEditor: false, Creator: testTicketCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Processing}

	ticket := ticket.CreateTestTicket(testTicketInfo, []ticket.MessageEntry{testMessage})

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicket", title, testTicketCreator, testMessage).Return(ticket, nil)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", mail).Return(false, -1)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, 302, resp.StatusCode, "Should return status code 302")
	assert.Equal(t, "/ticket/26", newLocation, "Should be redirected to \"/ticket/26\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Only POST requests should be possible
*/
func TestTicketCreateHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/create_ticket", nil)

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, -1)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	assert.Equal(t, 405, resp.StatusCode, "Should return status code 405")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Internal Only Parsing Error should return a 500.
*/
func TestTicketCreateHandler_ServeHTTP_ParsingError(t *testing.T) {
	mail := "test@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := "asdasdasdas"

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/ticket_create", newLocation, "Should be redirected to \"/ticket_create\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Internal Only empty should set it to false.
*/
func TestTicketCreateHandler_ServeHTTP_InternalOnlyEmpty(t *testing.T) {
	mail := "test@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := ""

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)

	testUser := user.User{Mail: mail, UserId: 5, FirstName: "Max", LastName: "Muller", Role: user.RegisteredUser, State: user.Active}
	testMessage := ticket.MessageEntry{Id: 0, CreatorMail: mail, Content: message, OnlyInternal: false}

	testTicketCreator := ticket.Creator{Mail: mail, FirstName: "Max", LastName: "Muller"}

	testTicketInfo := ticket.TicketInfo{Id: 26, Title: title, Editor: testUser, HasEditor: true, Creator: testTicketCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Processing}

	ticket := ticket.CreateTestTicket(testTicketInfo, []ticket.MessageEntry{testMessage})

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicketForInternalUser", title, testUser, testMessage).Return(ticket, nil)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", mail).Return(true, testUser.UserId)
	mockedUserContext.On("GetUserById", testUser.UserId).Return(true, testUser)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, 302, resp.StatusCode, "Should return status code 302")
	assert.Equal(t, "/ticket/26", newLocation, "Should be redirected to \"/ticket/26\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Create New Ticket throws an error should result in a 500.
*/
func TestTicketCreateHandler_ServeHTTP_CreateNewTicketError500(t *testing.T) {
	mail := "test@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := "false"
	firstName := "Max"
	lastName := "Muller"

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)
	req.Form.Add("first_name", firstName)
	req.Form.Add("last_name", lastName)

	testMessage := ticket.MessageEntry{Id: 0, CreatorMail: mail, Content: message, OnlyInternal: false}

	testTicketCreator := ticket.Creator{Mail: mail, FirstName: firstName, LastName: lastName}

	testTicketInfo := ticket.TicketInfo{Id: 26, Title: title, Editor: user.GetInvalidDefaultUser(), HasEditor: false, Creator: testTicketCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Processing}

	ticket := ticket.CreateTestTicket(testTicketInfo, []ticket.MessageEntry{testMessage})

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicket", title, testTicketCreator, testMessage).Return(ticket, errors.New("TestError"))
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", mail).Return(false, -1)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, 500, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/ticket_create", newLocation, "Should be redirected to \"/ticket_create\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	If a user tries to create a ticket for another user it should result in a 400.
*/
func TestTicketCreateHandler_ServeHTTP_LoggedInUserTriesToCreateTicketForOtherMail(t *testing.T) {
	mail := "test2@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := "false"

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", mail).Return(true, 4)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return status code 400")
	assert.Equal(t, "/ticket_create", newLocation, "Should be redirected to \"/ticket_create\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Create New Ticket For Internal User throws an error should result in a 500.
*/
func TestTicketCreateHandler_ServeHTTP_UserNotFoundWithId(t *testing.T) {
	mail := "test@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := "false"

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)

	testUser := user.User{Mail: mail, UserId: 5, FirstName: "Max", LastName: "Muller", Role: user.RegisteredUser, State: user.Active}

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", mail).Return(true, testUser.UserId)
	mockedUserContext.On("GetUserById", testUser.UserId).Return(false, *new(user.User))

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/ticket_create", newLocation, "Should be redirected to \"/ticket_create\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A Valid Http Request for logged In user.
*/
func TestTicketCreateHandler_ServeHTTP_CreateNewTicketForInternalUserError500(t *testing.T) {
	mail := "test@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := "false"

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)

	testUser := user.User{Mail: mail, UserId: 5, FirstName: "Max", LastName: "Muller", Role: user.RegisteredUser, State: user.Active}
	testMessage := ticket.MessageEntry{Id: 0, CreatorMail: mail, Content: message, OnlyInternal: false}

	testTicketCreator := ticket.Creator{Mail: mail, FirstName: "Max", LastName: "Muller"}

	testTicketInfo := ticket.TicketInfo{Id: 26, Title: title, Editor: testUser, HasEditor: true, Creator: testTicketCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Processing}

	ticket := ticket.CreateTestTicket(testTicketInfo, []ticket.MessageEntry{testMessage})

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicketForInternalUser", title, testUser, testMessage).Return(ticket, errors.New("TestError"))
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", mail).Return(true, testUser.UserId)
	mockedUserContext.On("GetUserById", testUser.UserId).Return(true, testUser)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, 500, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/ticket_create", newLocation, "Should be redirected to \"/ticket_create\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	If a not logged in user tries to create a ticket for another user it should result in a 400.
*/
func TestTicketCreateHandler_ServeHTTP_NotLoggedInUserTriesToCreateTicketForOtherMail(t *testing.T) {
	mail := "test2@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := "false"

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", mail).Return(true, 5)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return status code 400")
	assert.Equal(t, "/ticket_create", newLocation, "Should be redirected to \"/ticket_create\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}