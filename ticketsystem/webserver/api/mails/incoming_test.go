package mails

import (
	"bytes"
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"html"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockedMailFilter struct {
	mock.Mock
}

func (m *MockedMailFilter) IsAutomaticResponse(mail mailData.Mail) bool {
	args := m.Called(mail)
	return args.Bool(0)
}

/*
	Get a test handler with mocked data.
*/
func getTestHandlerWithMockedData() IncomingMailHandler {
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedMailFilter := new(MockedMailFilter)
	mockedMailFilter.On("IsAutomaticResponse", mock.Anything).Return(false)
	return IncomingMailHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext, TicketContext: mockedTicketContext,
		UserContext: mockedUserContext, MailRepliesFilter: mockedMailFilter}
}

/*
	Handling a mail to a existing ticketData should be able.
*/
func TestIncomingMailHandler_handleExistingTicketMail(t *testing.T) {
	testee := getTestHandlerWithMockedData()

	// Overwrite the mocked interface which are needed in this test:
	ticketId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticketData.Ticket{}, nil)

	testee.TicketContext = mockedTicketContext

	testMail := mailData.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}

	err := testee.handleExistingTicketMail(ticketId, testMail)
	assert.Nil(t, err)

	mockedTicketContext.AssertExpectations(t)
	// Assert that the parameter has been correctly set:
	assert.Equal(t, ticketId, mockedTicketContext.Calls[0].Arguments[0], "The correct ticketData id should be provided")

	assertMessageEntryIgnoringTime(t, testee.buildMessageEntry(testMail),
		mockedTicketContext.Calls[0].Arguments[1].(ticketData.MessageEntry))
}

/*
	Handling a mail from a registered userData should be possible.
*/
func TestIncomingMailHandler_handleNewTicketMail_RegisteredUser(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	userId := 1
	testMail := mailData.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}

	// Overwrite the mocked interface which are needed in this test:
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", testMail.Sender).Return(true, userId)
	returnedUser := userData.User{}
	mockedUserContext.On("GetUserById", userId).Return(true, returnedUser)

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicketForInternalUser", testMail.Subject, returnedUser, mock.Anything).
		Return(&ticketData.Ticket{}, nil)

	testee.TicketContext = mockedTicketContext
	testee.UserContext = mockedUserContext

	// Execute the test:
	err := testee.handleNewTicketMail(testMail)
	assert.Nil(t, err)

	// Assert the correct parameters
	mockedUserContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Handling a mail from a sender which is no registered userData should be able.
*/
func TestIncomingMailHandler_handleNewTicketMail_SenderIsNoUser(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	testMail := mailData.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}
	creator := testee.buildCreator(testMail)

	// Overwrite the mocked interface which are needed in this test:
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", testMail.Sender).Return(false, -1)

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicket", testMail.Subject, creator, mock.Anything).
		Return(&ticketData.Ticket{}, nil)

	// Set the mocked contexts:
	testee.TicketContext = mockedTicketContext
	testee.UserContext = mockedUserContext

	// Execute the test:
	err := testee.handleNewTicketMail(testMail)
	assert.Nil(t, err)

	// Assert the correct parameters
	mockedUserContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Handling a mail from a userData which is registered, should be possible.
*/
func TestIncomingMailHandler_handleIncomingMails_SenderIsAUser(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	userId := 1
	testMail := mailData.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}
	testMails := []mailData.Mail{testMail}

	// Overwrite the mocked interface which are needed in this test:
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", testMail.Sender).Return(true, userId)
	returnedUser := userData.User{}
	mockedUserContext.On("GetUserById", userId).Return(true, returnedUser)

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicketForInternalUser", testMail.Subject, returnedUser, mock.Anything).
		Return(&ticketData.Ticket{}, nil)

	testee.TicketContext = mockedTicketContext
	testee.UserContext = mockedUserContext

	// Execute the test:
	err := testee.handleIncomingMails(testMails)
	assert.Nil(t, err)

	// Assert the correct parameters
	mockedUserContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	Handling a mail for when the sender is no userData, should be possible.
*/
func TestIncomingMailHandler_handleIncomingMails_SenderIsNoUser(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	testMail := mailData.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}
	creator := testee.buildCreator(testMail)
	testMails := []mailData.Mail{testMail}

	// Overwrite the mocked interface which are needed in this test:
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", testMail.Sender).Return(false, -1)

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicket", testMail.Subject, creator, mock.Anything).
		Return(&ticketData.Ticket{}, nil)

	// Set the mocked contexts:
	testee.TicketContext = mockedTicketContext
	testee.UserContext = mockedUserContext

	// Execute the test:
	err := testee.handleIncomingMails(testMails)
	assert.Nil(t, err)

	// Assert the correct parameters
	mockedUserContext.AssertExpectations(t)
	mockedTicketContext.AssertExpectations(t)
}

/*
	A incoming mail for a existing ticketData should notify the creator of the ticketData.
*/
func TestIncomingMailHandler_handleIncomingMails_TicketExists_MailSent(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	testMail := mailData.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "Ticket<1> TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}
	testMails := []mailData.Mail{testMail}
	existingTicket := ticketData.Ticket{}

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("GetTicketById", 1).Return(true, &existingTicket)
	mockedTicketContext.On("AppendMessageToTicket", 1, mock.Anything).
		Return(&existingTicket, nil)

	mockedMailContext := new(mockedForTests.MockedMailContext)
	expectedSubject := "New Entry for your ticketData: " + html.EscapeString(testMail.Subject)
	expectedMailContent := mailData.BuildAppendMessageNotificationMailContent(existingTicket.Info().Creator.Mail, testMail.Sender, testMail.Content)

	mockedMailContext.On("CreateNewOutgoingMail", "", expectedSubject, expectedMailContent).Return(nil)

	testee.TicketContext = mockedTicketContext
	testee.MailContext = mockedMailContext

	// Execute the test:
	err := testee.handleIncomingMails(testMails)
	assert.Nil(t, err)

	// Assert the correct parameters
	mockedTicketContext.AssertExpectations(t)
	mockedMailContext.AssertExpectations(t)
}

/*
	A invalid payload should result in a 400.
*/
func TestIncomingMailHandler_ServeHTTP_InvalidPayload_400Returned(t *testing.T) {
	testee := getTestHandlerWithMockedData()

	jsonData, _ := json.Marshal(getTestMails())
	// Make the jsonData invalid:
	jsonData[2] = 4

	req, err := http.NewRequest("POST", shared.SendPath, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 400, rr.Code, "Status code 400 should be returned")
}

/*
	A error should result in a 500.
*/
func TestIncomingMailHandler_ServeHTTP_500Returned(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("GetTicketById", mock.Anything).Return(true, &ticketData.Ticket{})
	mockedTicketContext.On("AppendMessageToTicket", mock.Anything, mock.Anything).
		Return(&ticketData.Ticket{}, errors.New("TestError"))
	testee.TicketContext = mockedTicketContext

	jsonData, _ := json.Marshal(getTestMails())

	req, err := http.NewRequest("POST", shared.SendPath, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")
}

/*
	The filter should be respected. The other mocks should not be called.
*/
func TestIncomingMailHandler_ServeHTTP_MailFiltered(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	mockedFilter := new(MockedMailFilter)
	mockedFilter.On("IsAutomaticResponse", mock.Anything).Return(true)
	testee.MailRepliesFilter = mockedFilter

	jsonData, _ := json.Marshal(getTestMails())

	req, err := http.NewRequest("POST", shared.SendPath, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedFilter.AssertExpectations(t)
}

/*
	Assert a MessageEntry ignoring the autonatically set timestamp.
*/
func assertMessageEntryIgnoringTime(t *testing.T, expected ticketData.MessageEntry, actual ticketData.MessageEntry) {
	assert.Equal(t, expected.Content, actual.Content, "Content should be equal")
	assert.Equal(t, expected.CreatorMail, actual.CreatorMail, "CreatorMail should be equal")
	assert.Equal(t, expected.OnlyInternal, actual.OnlyInternal, "OnlyInternal should be equal")
	assert.Equal(t, expected.CreatorMail, actual.CreatorMail, "CreatorMail should be equal")
}
