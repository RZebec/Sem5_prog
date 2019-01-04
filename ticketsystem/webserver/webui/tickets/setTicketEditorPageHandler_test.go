package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

/*
	Set a editor.
*/
func TestTicketSetEditorHandler_ServeHTTP_SetEditor_ValidRequest(t *testing.T) {
	ticketId := 4
	editorUserId := 8
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", editorUserId).Return(true, user.User{UserId: editorUserId})
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})
	// A notification mail should be sent
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Two changes should be stored in the history:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil).Twice()
	// The editor should be set:
	mockedTicketContext.On("SetEditor", mock.Anything, ticketId).Return(&ticket.Ticket{}, nil)
	// The state of the ticket should be changed:
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Processing).Return(&ticket.Ticket{}, nil)

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 200, resp.StatusCode, "Should return 200")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)

	setEditorId := mockedTicketContext.Calls[1].Arguments[0].(user.User).UserId
	assert.Equal(t, editorUserId, setEditorId, "The editor should be set")
}

/*
	Remove a editor.
*/
func TestTicketSetEditorHandler_ServeHTTP_RemoveEditor_ValidRequest(t *testing.T) {
	ticketId := 4
	// Setting editor to -1 should remove it.
	editorUserId := -1
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})
	// A notification mail should be sent
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	// Two changes should be stored in the history:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil).Twice()
	// The editor should be set:
	mockedTicketContext.On("RemoveEditor", ticketId).Return(nil)
	// The state of the ticket should be changed:
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Processing).Return(&ticket.Ticket{}, nil)

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 200, resp.StatusCode, "Should return 200")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Setting a non existing user should not be possible.
*/
func TestTicketSetEditorHandler_ServeHTTP_UnknownEditorId_InvalidRequest(t *testing.T) {
	ticketId := 4
	editorUserId := 8
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", editorUserId).Return(false, user.User{})

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Execution with a non existing logged in user should not be possible. This combination should not be possible.
*/
func TestTicketSetEditorHandler_ServeHTTP_UnknownLoggedInUser_invalidRequest(t *testing.T) {
	ticketId := 4
	editorUserId := 8
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", editorUserId).Return(true, user.User{UserId: editorUserId})
	mockedUserContext.On("GetUserById", loggedInUserId).Return(false, user.User{})

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Error in the context should result in a 500.
*/
func TestTicketSetEditorHandler_ServeHTTP_SetEditorFails_Returns500(t *testing.T) {
	ticketId := 4
	editorUserId := 8
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", editorUserId).Return(true, user.User{UserId: editorUserId})
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})
	mockedTicketContext.On("SetEditor", mock.Anything, ticketId).Return(&ticket.Ticket{}, errors.New("TestError"))

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Error during message appending should result in a 500
*/
func TestTicketSetEditorHandler_ServeHTTP_AppendMessageFailed_500Returned(t *testing.T) {
	ticketId := 4
	editorUserId := 8
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", editorUserId).Return(true, user.User{UserId: editorUserId})
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})

	// Appending to the history should fail:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).
		Return(&ticket.Ticket{}, errors.New("Testerror"))
	// The editor should be set:
	mockedTicketContext.On("SetEditor", mock.Anything, ticketId).Return(&ticket.Ticket{}, nil)

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	A error during the change of the ticket state should result in a 500.
*/
func TestTicketSetEditorHandler_ServeHTTP_TicketStateChangedFailed_Returns500(t *testing.T) {
	ticketId := 4
	editorUserId := 8
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", editorUserId).Return(true, user.User{UserId: editorUserId})
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})

	// Two changes should be stored in the history:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil).Once()
	// The editor should be set:
	mockedTicketContext.On("SetEditor", mock.Anything, ticketId).Return(&ticket.Ticket{}, nil)

	// Ticket state change should fail
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Processing).
		Return(&ticket.Ticket{}, errors.New("TestError"))

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	A error during the change of the ticket state history writing should result in a 500.
*/
func TestTicketSetEditorHandler_ServeHTTP_TicketStateChangedHistoryWriteFailed_Returns500(t *testing.T) {
	ticketId := 4
	editorUserId := 8
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", editorUserId).Return(true, user.User{UserId: editorUserId})
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})

	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil).Once()
	// The editor should be set:
	mockedTicketContext.On("SetEditor", mock.Anything, ticketId).Return(&ticket.Ticket{}, nil)

	// Ticket state change should work
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Processing).
		Return(&ticket.Ticket{}, nil)

	// Second message for the history should fail
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).
		Return(&ticket.Ticket{}, errors.New("TestError")).Once()

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	A error during the notification should result in a 500.
*/
func TestTicketSetEditorHandler_ServeHTTP_NotificationFailed_Returns500(t *testing.T) {
	ticketId := 4
	editorUserId := 8
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", editorUserId).Return(true, user.User{UserId: editorUserId})
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})
	// Two changes should be stored in the history:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil).Twice()
	// The editor should be set:
	mockedTicketContext.On("SetEditor", mock.Anything, ticketId).Return(&ticket.Ticket{}, nil)
	// The state of the ticket should be changed:
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Processing).Return(&ticket.Ticket{}, nil)

	// A notification mail should fail
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("TestError"))

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)

	setEditorId := mockedTicketContext.Calls[1].Arguments[0].(user.User).UserId
	assert.Equal(t, editorUserId, setEditorId, "The editor should be set")
}

/*
	Error during the editor removing should result in a 500.
*/
func TestTicketSetEditorHandler_ServeHTTP_RemoveEditorFails_Returns500(t *testing.T) {
	ticketId := 4
	// Setting editor to -1 should remove it.
	editorUserId := -1
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	// Removing the editor should fail:
	mockedTicketContext.On("RemoveEditor", ticketId).Return(errors.New("TestError"))

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Execution for a logged in user which does not exist should not be possible. This should never occur.
*/
func TestTicketSetEditorHandler_ServeHTTP_RemoveEditorFails_InvalidRequest(t *testing.T) {
	ticketId := 4
	// Setting editor to -1 should remove it.
	editorUserId := -1
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	// Removing the editor should work:
	mockedTicketContext.On("RemoveEditor", ticketId).Return(nil)
	// The user should not exist:
	mockedUserContext.On("GetUserById", loggedInUserId).Return(false, user.User{})

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Error during message appending should result in a 500.
*/
func TestTicketSetEditorHandler_ServeHTTP_RemoveEditorFails_500Returned(t *testing.T) {
	ticketId := 4
	// Setting editor to -1 should remove it.
	editorUserId := -1
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	// Removing the editor should work:
	mockedTicketContext.On("RemoveEditor", ticketId).Return(nil)
	// The user should exist:
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})
	// Appending the message should fail:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).
		Return(&ticket.Ticket{}, errors.New("TestError"))

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Error during state change should result in a 500.
*/
func TestTicketSetEditorHandler_ServeHTTP_StateChangeFails_500Returned(t *testing.T) {
	ticketId := 4
	// Setting editor to -1 should remove it.
	editorUserId := -1
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	// Removing the editor should work:
	mockedTicketContext.On("RemoveEditor", ticketId).Return(nil)
	// The user should exist:
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})
	// Appending the message should work:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).
		Return(&ticket.Ticket{}, nil)
	// Changing the state should fail:
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Processing).
		Return(&ticket.Ticket{}, errors.New("TestEditor"))

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Error during message appending of state change should result in a 500.
*/
func TestTicketSetEditorHandler_ServeHTTP_StateChangeHistoryFails_500Returned(t *testing.T) {
	ticketId := 4
	// Setting editor to -1 should remove it.
	editorUserId := -1
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	// Removing the editor should work:
	mockedTicketContext.On("RemoveEditor", ticketId).Return(nil)
	// The user should exist:
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})
	// Appending the message should work:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).
		Return(&ticket.Ticket{}, nil).Once()
	// Changing the state should work:
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Processing).
		Return(&ticket.Ticket{}, nil)
	// Appending the second message should fail:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).
		Return(&ticket.Ticket{}, errors.New("TestError"))

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	A error during the mail notification should result in a 500.
*/
func TestTicketSetEditorHandler_ServeHTTP_NotificationFailed_500Returned(t *testing.T) {
	ticketId := 4
	// Setting editor to -1 should remove it.
	editorUserId := -1
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{})
	// Two changes should be stored in the history:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil).Twice()
	// The editor should be set:
	mockedTicketContext.On("RemoveEditor", ticketId).Return(nil)
	// The state of the ticket should be changed:
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Processing).Return(&ticket.Ticket{}, nil)

	// A notification mail should fail:
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("TestError"))

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))
}

/*
	A invalid ticket id should result in a invalid request.
*/
func TestTicketSetEditorHandler_ServeHTTP_InvalidTicketId_InvalidRequest(t *testing.T) {
	editorUserId := 8
	loggedInUserId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "asd")
	req.Form.Add("editorUserId", strconv.Itoa(editorUserId))
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	A invalid user id should result in a invalid request.
*/
func TestTicketSetEditorHandler_ServeHTTP_InvalidUserId_InvalidRequest(t *testing.T) {
	loggedInUserId := 2
	ticketId := 3
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", "sfda")
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	A ticket id for a non existing ticket should result in a invalid request.
*/
func TestTicketSetEditorHandler_ServeHTTP_NonExistingTicket_InvalidRequest(t *testing.T) {
	loggedInUserId := 2
	ticketId := 3
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	// The ticket should not exist:
	mockedTicketContext.On("GetTicketById", ticketId).Return(false, &ticket.Ticket{})

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("editorUserId", "sfda")
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Only post request should be possible.
*/
func TestTicketSetEditorHandler_ServeHTTP_InvalidRequestMethod_InvalidRequest(t *testing.T) {

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketSetEditorHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("Get", "/test", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Should return 405")

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}
