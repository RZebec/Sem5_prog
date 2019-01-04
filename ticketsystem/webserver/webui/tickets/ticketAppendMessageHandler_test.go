package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"

	"strconv"
	"testing"
)

/*
	A valid request from a authenticated user should append the message to the ticket.
*/
func TestTicketAppendMessageHandler_ServeHTTP_AuthenticatedUser_ValidRequest(t *testing.T) {
	ticketId := 2
	testMessageContent := "TestMessageEntryContent"
	onlyInternal := true
	userId := 5
	userMail := "test@test.de"
	userFromContext := user.User{UserId: 5, Mail: userMail}

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", userId).Return(true, userFromContext)
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", testMessageContent)
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))
	// Add authentication info:
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, userId)

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 200, resp.StatusCode, "Should return 200")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)

	// Assert that the message entry has been created:
	assert.Equal(t, ticketId, mockedTicketContext.Calls[1].Arguments[0].(int), "The correct ticket id should be used.")
	messageEntry := mockedTicketContext.Calls[1].Arguments[1].(ticket.MessageEntry)
	assert.Equal(t, userMail, messageEntry.CreatorMail, "The correct user mail should be set")
	assert.Equal(t, testMessageContent, messageEntry.Content, "The correct content should be set")
	assert.Equal(t, onlyInternal, messageEntry.OnlyInternal, "The correct onlyInternal flag should be set")
	assert.NotNil(t, messageEntry.CreationTime, "The creation time should be set")
}

/*
	A error during the execution should result in a 500.
*/
func TestTicketAppendMessageHandler_ServeHTTP_AuthenticatedUser_ContextReturnsError_Returns500(t *testing.T) {
	ticketId := 2
	testMessageContent := "TestMessageEntryContent"
	onlyInternal := true
	userId := 5
	userMail := "test@test.de"
	userFromContext := user.User{UserId: 5, Mail: userMail}

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", userId).Return(true, userFromContext)
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{},
		errors.New("TestError"))

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", testMessageContent)
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))
	// Add authentication info:
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, userId)

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A request for a non existing ticket is invalid.
*/
func TestTicketAppendMessageHandler_ServeHTTP_AuthenticatedUser_TicketDoesNotExist_InvalidRequest(t *testing.T) {
	ticketId := 2
	testMessageContent := "TestMessageEntryContent"
	onlyInternal := true
	userId := 5

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(false, &ticket.Ticket{})

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", testMessageContent)
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))
	// Add authentication info:
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, userId)

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A request with an invalid ticket id is invalid.
*/
func TestTicketAppendMessageHandler_ServeHTTP_AuthenticatedUser_InvalidTicketId_InvalidRequest(t *testing.T) {
	testMessageContent := "TestMessageEntryContent"
	onlyInternal := true
	userId := 5

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "s2")
	req.Form.Add("messageContent", testMessageContent)
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))
	// Add authentication info:
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, userId)

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	If the user does not exist -> invalid request. Should not be possible.
*/
func TestTicketAppendMessageHandler_ServeHTTP_AuthenticatedUser_UserDoesNotExist_InvalidRequest(t *testing.T) {
	ticketId := 2
	testMessageContent := "TestMessageEntryContent"
	onlyInternal := true
	userId := 5
	userMail := "test@test.de"
	userFromContext := user.User{UserId: 5, Mail: userMail}

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", userId).Return(false, userFromContext)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", testMessageContent)
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))
	// Add authentication info:
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, userId)

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Appending a empty message is invalid.
*/
func TestTicketAppendMessageHandler_ServeHTTP_AuthenticatedUser_EmptyContent_InvalidRequest(t *testing.T) {
	ticketId := 2
	onlyInternal := true
	userId := 5
	userMail := "test@test.de"
	userFromContext := user.User{UserId: 5, Mail: userMail}

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", userId).Return(true, userFromContext)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", "")
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))
	// Add authentication info:
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, userId)

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A invalid OnlyInternalFlag results in a invalid request.
*/
func TestTicketAppendMessageHandler_ServeHTTP_AuthenticatedUser_InvalidOnlyInternalFlag_InvalidRequest(t *testing.T) {
	ticketId := 2
	userId := 5
	userMail := "test@test.de"
	userFromContext := user.User{UserId: 5, Mail: userMail}

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", userId).Return(true, userFromContext)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", "TestContent")
	req.Form.Add("onlyInternal", "test")
	// Add authentication info:
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, userId)

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A non authenticated user should be possible to append a message.
*/
func TestTicketAppendMessageHandler_ServeHTTP_NonAuthenticatedUser_ValidRequest(t *testing.T) {
	ticketId := 2
	testMessageContent := "TestMessageEntryContent"
	userId := 5
	userMail := "test@test.de"

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil)
	mockedUserContext.On("GetUserForEmail", userMail).Return(false, userId)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", testMessageContent)
	req.Form.Add("mail", userMail)

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, 200, resp.StatusCode, "Should return 200")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)

	// Assert that the message entry has been created:
	assert.Equal(t, ticketId, mockedTicketContext.Calls[1].Arguments[0].(int), "The correct ticket id should be used.")
	messageEntry := mockedTicketContext.Calls[1].Arguments[1].(ticket.MessageEntry)
	assert.Equal(t, userMail, messageEntry.CreatorMail, "The correct user mail should be set")
	assert.Equal(t, testMessageContent, messageEntry.Content, "The correct content should be set")
	assert.NotNil(t, messageEntry.CreationTime, "The creation time should be set")
}

/*
	A error from the context should result in a 500.
*/
func TestTicketAppendMessageHandler_ServeHTTP_NonAuthenticatedUser_ContextReturnsError_Return500(t *testing.T) {
	ticketId := 2
	testMessageContent := "TestMessageEntryContent"
	onlyInternal := true
	userId := 5
	userMail := "test@test.de"

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserForEmail", userMail).Return(false, userId)
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{},
		errors.New("TestError"))

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("mail", userMail)
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", testMessageContent)
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A request for a ticket which does not exist is invalid.
*/
func TestTicketAppendMessageHandler_ServeHTTP_NonAuthenticatedUser_TicketDoesNotExist_InvalidRequest(t *testing.T) {
	ticketId := 2
	testMessageContent := "TestMessageEntryContent"
	onlyInternal := true

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(false, &ticket.Ticket{})

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", testMessageContent)
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A invalid ticket id results in a invalid request.
*/
func TestTicketAppendMessageHandler_ServeHTTP_NonAuthenticatedUser_InvalidTicketId_InvalidRequest(t *testing.T) {
	testMessageContent := "TestMessageEntryContent"
	onlyInternal := true

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "s2")
	req.Form.Add("messageContent", testMessageContent)
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A empty message should not be appended.
*/
func TestTicketAppendMessageHandler_ServeHTTP_NonAuthenticatedUser_EmptyContent_InvalidRequest(t *testing.T) {
	ticketId := 2
	onlyInternal := true
	userId := 5
	userMail := "test@test.de"

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserForEmail", userMail).Return(false, userId)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("mail", userMail)
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", "")
	req.Form.Add("onlyInternal", strconv.FormatBool(onlyInternal))

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	The only internal flag should be ignored for non authenticated users.
*/
func TestTicketAppendMessageHandler_ServeHTTP_NonAuthenticatedUser_OnlyInternalSet_FlagIgnored(t *testing.T) {
	ticketId := 2
	userId := 5
	userMail := "test@test.de"

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserForEmail", userMail).Return(false, userId)
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("mail", userMail)
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", "TestContent")
	req.Form.Add("onlyInternal", "true")

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, 200, resp.StatusCode, "Should return 200")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)

	// Assert that the message entry has been created:
	assert.Equal(t, ticketId, mockedTicketContext.Calls[1].Arguments[0].(int), "The correct ticket id should be used.")
	messageEntry := mockedTicketContext.Calls[1].Arguments[1].(ticket.MessageEntry)
	assert.Equal(t, userMail, messageEntry.CreatorMail, "The correct user mail should be set")
	assert.Equal(t, "TestContent", messageEntry.Content, "The correct content should be set")
	assert.Equal(t, false, messageEntry.OnlyInternal, "OnlyInternal flag should be set to false")
	assert.NotNil(t, messageEntry.CreationTime, "The creation time should be set")
}

/*
	A invalid mail should result in a invalid request.
*/
func TestTicketAppendMessageHandler_ServeHTTP_NonAuthenticatedUser_MailIsInvalid_InvalidRequest(t *testing.T) {
	ticketId := 2
	userMail := "@test@test.de"

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("mail", userMail)
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", "TestContent")
	req.Form.Add("onlyInternal", "true")

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/ticket/2", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Trying to use the mail of a authenticated user, while the current user is not authenticated, should result in
	redirect to the login page.
*/
func TestTicketAppendMessageHandler_ServeHTTP_NonAuthenticatedWithMailOfExistingUSer_RedirectedToLogin(t *testing.T) {
	ticketId := 2
	userId := 5
	userMail := "test@test.de"

	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserForEmail", userMail).Return(true, userId)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("mail", userMail)
	req.Form.Add("ticketId", "2")
	req.Form.Add("messageContent", "TestContent")
	req.Form.Add("onlyInternal", "true")

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, 200, resp.StatusCode, "Should return 200")
	assert.Equal(t, "/login", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Only post should be possible.
*/
func TestTicketAppendMessageHandler_ServeHTTP_InvalidRequestMethod(t *testing.T) {
	// Create and setup mocked interfaces:
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketAppendMessageHandler{UserContext: mockedUserContext, TicketContext: mockedTicketContext,
		Logger: testhelpers.GetTestLogger()}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	req, err := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err)

	// Execute the test and assert the result
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, 405, resp.StatusCode, "Should return 405")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}
