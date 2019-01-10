// 5894619, 6720876, 9793350
package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
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
	Tickets should be able to be merged.
*/
func TestTickerMergeHandler_ServeHTTP_TicketsMerged(t *testing.T) {
	firstTicketId := 22
	secondTicketId := 4
	loggedInUserId := 9
	loggedInUserMail := "1234test@test.de"

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).Return(true, nil)
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, userData.User{Mail: loggedInUserMail})
	mockedTicketContext.On("AppendMessageToTicket", secondTicketId, mock.Anything).Return(&ticketData.Ticket{}, nil)

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")
	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode, "Should return 302")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	// Assert the appended message:
	appendedMessage := mockedTicketContext.Calls[3].Arguments[1].(ticketData.MessageEntry)
	assert.Contains(t, appendedMessage.Content, "Tickets merged", "A ticket merged message should be appended")
	assert.Contains(t, appendedMessage.Content, "22", "Id of the first ticket should be in the message")
	assert.Contains(t, appendedMessage.Content, "4", "Id of the second ticket should be in the message")

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Invalid ticket id should result in a invalid request.
*/
func TestTickerMergeHandler_ServeHTTP_InvalidFirstTicketId_InvalidRequest(t *testing.T) {
	loggedInUserId := 9
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", "sdsd")
	req.Form.Add("secondTicketId", "3")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Invalid ticket id should result in a invalid request.
*/
func TestTickerMergeHandler_ServeHTTP_InvalidSecondTicketId_InvalidRequest(t *testing.T) {
	loggedInUserId := 9
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", "2")
	req.Form.Add("secondTicketId", "")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A id of a non existing ticket should result in a invalid request.
*/
func TestTickerMergeHandler_ServeHTTP_FirstTicketDoesNotExist_InvalidRequest(t *testing.T) {
	loggedInUserId := 9
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(false, &ticketData.Ticket{})

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A id of a non existing ticket should result in a invalid request.
*/
func TestTickerMergeHandler_ServeHTTP_SecondTicketDoesNotExist_InvalidRequest(t *testing.T) {
	loggedInUserId := 9
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(false, &ticketData.Ticket{})

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A error during the merge should result in a 500.
*/
func TestTickerMergeHandler_ServeHTTP_ContextReturnsError_Returns500(t *testing.T) {
	loggedInUserId := 9
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).
		Return(false, errors.New("TestError"))

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Error during mail notification should result in a 500.
*/
func TestTickerMergeHandler_ServeHTTP_FirstMailNotSent_Returns500(t *testing.T) {
	loggedInUserId := 9
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).
		Return(true, nil)
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("TestError"))

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Error during mail notification should result in a 500.
*/
func TestTickerMergeHandler_ServeHTTP_SecondMailNotSent_Returns500(t *testing.T) {
	loggedInUserId := 9
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).
		Return(true, nil)
	// First mail is successfull:
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	// Second mail is not successfull:
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("TestError")).Once()

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	Only post method should be possible.
*/
func TestTickerMergeHandler_ServeHTTP_InvalidRequestMethod(t *testing.T) {
	loggedInUserId := 9
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")

	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Should return method not allowed")

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A user which is logged in but does not exist, should result in a invalid request. This should never be possible.
*/
func TestTickerMergeHandler_ServeHTTP_InvalidLoggedInUser_InvalidRequest(t *testing.T) {
	firstTicketId := 22
	secondTicketId := 4
	loggedInUserId := 9
	loggedInUserMail := "1234test@test.de"

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).Return(true, nil)
	// Return false, so the user does not exist:s
	mockedUserContext.On("GetUserById", loggedInUserId).Return(false, userData.User{Mail: loggedInUserMail})

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")
	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}

/*
	A error during the appending of the history message, should result in a 500.
*/
func TestTickerMergeHandler_ServeHTTP_ErrorDuringMessageAppending_Returns500(t *testing.T) {
	firstTicketId := 22
	secondTicketId := 4
	loggedInUserId := 9
	loggedInUserMail := "1234test@test.de"

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)

	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).Return(true, nil)
	mockedUserContext.On("GetUserById", loggedInUserId).Return(true, userData.User{Mail: loggedInUserMail})
	// Append throws an error:
	mockedTicketContext.On("AppendMessageToTicket", secondTicketId, mock.Anything).
		Return(&ticketData.Ticket{}, errors.New("TestError"))

	testee := TicketMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext, UserContext: mockedUserContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, loggedInUserId, "")
	// Execute the test:
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)
}
