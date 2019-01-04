package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
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

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)

	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).Return(true, nil)

	testee := TickerMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode, "Should return 302")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}

/*
	Invalid ticket id should result in a invalid request.
*/
func TestTickerMergeHandler_ServeHTTP_InvalidFirstTicketId_InvalidRequest(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)

	testee := TickerMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", "sdsd")
	req.Form.Add("secondTicketId", "3")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}

/*
	Invalid ticket id should result in a invalid request.
*/
func TestTickerMergeHandler_ServeHTTP_InvalidSecondTicketId_InvalidRequest(t *testing.T) {
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)

	testee := TickerMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", "2")
	req.Form.Add("secondTicketId", "")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}

/*
	A id of a non existing ticket should result in a invalid request.
*/
func TestTickerMergeHandler_ServeHTTP_FirstTicketDoesNotExist_InvalidRequest(t *testing.T) {
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(false, &ticket.Ticket{})

	testee := TickerMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}

/*
	A id of a non existing ticket should result in a invalid request.
*/
func TestTickerMergeHandler_ServeHTTP_SecondTicketDoesNotExist_InvalidRequest(t *testing.T) {
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(false, &ticket.Ticket{})

	testee := TickerMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}

/*
	A error during the merge should result in a 500.
*/
func TestTickerMergeHandler_ServeHTTP_ContextReturnsError_Returns500(t *testing.T) {
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).
		Return(false, errors.New("TestError"))

	testee := TickerMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}

/*
	Error during mail notification should result in a 500.
*/
func TestTickerMergeHandler_ServeHTTP_FirstMailNotSent_Returns500(t *testing.T) {
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).
		Return(true, nil)
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("TestError"))

	testee := TickerMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}

/*
	Error during mail notification should result in a 500.
*/
func TestTickerMergeHandler_ServeHTTP_SecondMailNotSent_Returns500(t *testing.T) {
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)

	mockedTicketContext.On("GetTicketById", firstTicketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("GetTicketById", secondTicketId).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("MergeTickets", firstTicketId, secondTicketId).
		Return(true, nil)
	// First mail is successfull:
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	// Second mail is not successfull:
	mockedMailContext.On("CreateNewOutgoingMail", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("TestError")).Once()

	testee := TickerMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext}

	req, err := http.NewRequest("POST", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, 500, resp.StatusCode, "Should return 500")
	assert.Equal(t, "/ticket/4", resp.Header.Get("location"))

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}

/*
	Only post method should be possible.
*/
func TestTickerMergeHandler_ServeHTTP_InvalidRequestMethod(t *testing.T) {
	firstTicketId := 22
	secondTicketId := 4

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedMailContext := new(mockedForTests.MockedMailContext)

	testee := TickerMergeHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext,
		TicketContext: mockedTicketContext}

	req, err := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err)

	// Add the values to the request:
	req.Form = url.Values{}
	req.Form.Add("firstTicketId", strconv.Itoa(firstTicketId))
	req.Form.Add("secondTicketId", strconv.Itoa(secondTicketId))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	// Execute the test:
	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Should return method not allowed")

	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}
