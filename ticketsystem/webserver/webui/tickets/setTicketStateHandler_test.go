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
	Changing the state of the ticket should be possible.
*/
func TestSetTicketStateHandler_ServeHTTP_ValidStateSet(t *testing.T) {
	ticketId := 5
	loggedInUserId := 3
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	// The ticket exists:
	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	// The user exists:
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{UserId: loggedInUserId})
	// Changing the ticket state should be successfull:
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Closed).Return(&ticket.Ticket{}, nil)
	// The history message should be appended:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil)

	testee := SetTicketStateHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("newState", "Closed")
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 200, resp.StatusCode, "Should return 200")
	assert.Equal(t, "/ticket/5", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Changing the state to a invalid state should not be possbile.
*/
func TestSetTicketStateHandler_ServeHTTP_InvalidState_InvalidRequest(t *testing.T) {
	ticketId := 5
	loggedInUserId := 3
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	// The ticket exists:
	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	// The user exists:
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{UserId: loggedInUserId})

	testee := SetTicketStateHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("newState", "perfect")
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/ticket/5", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	A invalid ticket id should result in a invalid request.
*/
func TestSetTicketStateHandler_ServeHTTP_InvalidTicketId_InvalidRequest(t *testing.T) {
	loggedInUserId := 3
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := SetTicketStateHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", "asd")
	req.Form.Add("newState", "perfect")
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
	Setting the state of a non existing ticket should not be possible.
*/
func TestSetTicketStateHandler_ServeHTTP_NonExistingTicket_InvalidRequest(t *testing.T) {
	ticketId := 5
	loggedInUserId := 3
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	// The ticket exists:
	mockedTicketContext.On("GetTicketById", ticketId).Return(false, &ticket.Ticket{})

	testee := SetTicketStateHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("newState", "Closed")
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
	A error during the state change should result in a 500.
*/
func TestSetTicketStateHandler_ServeHTTP_ErrorDuringStateChange_Returns500(t *testing.T) {
	ticketId := 5
	loggedInUserId := 3
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	// The ticket exists:
	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	// Error during state change:
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Processing).
		Return(&ticket.Ticket{}, errors.New("TestError"))

	// The user does exist:
	mockedUserContext.On("GetUserById", loggedInUserId).
		Return(true, user.User{UserId: loggedInUserId})

	testee := SetTicketStateHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("newState", "Processing")
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/5", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	A error during appending of history message should result in a 500.
*/
func TestSetTicketStateHandler_ServeHTTP_AppendingMessageFailed_Returns500(t *testing.T) {
	ticketId := 5
	loggedInUserId := 3
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	// The ticket exists:
	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})
	// The user exists:
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, user.User{UserId: loggedInUserId})
	// Changing the ticket state should be successfull:
	mockedTicketContext.On("SetTicketState", ticketId, ticket.Open).Return(&ticket.Ticket{}, nil)
	// Error during message appending:
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).
		Return(&ticket.Ticket{}, errors.New("TestError"))

	testee := SetTicketStateHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("newState", "open")
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/5", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Changing the state with a user which does not exist, should fail. This should never happen.
*/
func TestSetTicketStateHandler_ServeHTTP_UnknownLoggedInUser_InvalidRequest(t *testing.T) {
	ticketId := 5
	loggedInUserId := 3
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	// The ticket exists:
	mockedTicketContext.On("GetTicketById", ticketId).Return(true, &ticket.Ticket{})

	// The user does not exists:
	mockedUserContext.On("GetUserById", loggedInUserId).Return(false, user.User{UserId: loggedInUserId})

	testee := SetTicketStateHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("ticketId", strconv.Itoa(ticketId))
	req.Form.Add("newState", "Closed")
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/ticket/5", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Only post should be possible.
*/
func TestSetTicketStateHandler_InvalidRequestMethod_ValidStateSet(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := SetTicketStateHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("get", "/test", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Should return 405")
	assert.Equal(t, "", resp.Header.Get("location"))

	mockedUserContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}
