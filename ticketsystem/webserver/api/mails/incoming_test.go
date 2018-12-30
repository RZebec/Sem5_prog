package mails

import (
	"bytes"
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
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

/*
	Get a test handler with mocked data.
*/
func getTestHandlerWithMockedData() IncomingMailHandler {
	mockedMailContext := new(mockedForTests.MockedMailContext)
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	return IncomingMailHandler{Logger: getTestLogger(), MailContext: mockedMailContext, TicketContext: mockedTicketContext, UserContext: mockedUserContext}
}

/*
	Handling a mail to a existing ticket should be able.
*/
func TestIncomingMailHandler_handleExistingTicketMail(t *testing.T) {
	testee := getTestHandlerWithMockedData()

	// Overwrite the mocked interface which are needed in this test:
	ticketId := 2
	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("AppendMessageToTicket", ticketId, mock.Anything).Return(&ticket.Ticket{}, nil)

	testee.TicketContext = mockedTicketContext

	testMail := mail.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}

	err := testee.handleExistingTicketMail(ticketId, testMail)
	assert.Nil(t, err)

	mockedTicketContext.AssertExpectations(t)
	// Assert that the parameter has been correctly set:
	assert.Equal(t, ticketId, mockedTicketContext.Calls[0].Arguments[0], "The correct ticket id should be provided")

	assertMessageEntryIgnoringTime(t, testee.buildMessageEntry(testMail),
		mockedTicketContext.Calls[0].Arguments[1].(ticket.MessageEntry))
}

/*
	Handling a mail from a registered user should be possible.
*/
func TestIncomingMailHandler_handleNewTicketMail_RegisteredUser(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	userId := 1
	testMail := mail.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}

	// Overwrite the mocked interface which are needed in this test:
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", testMail.Sender).Return(true, userId)
	returnedUser := user.User{}
	mockedUserContext.On("GetUserById", userId).Return(true, returnedUser)

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicketForInternalUser", testMail.Subject, returnedUser, mock.Anything).
		Return(&ticket.Ticket{}, nil)

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
	Handling a mail from a sender which is no registered user should be able.
*/
func TestIncomingMailHandler_handleNewTicketMail_SenderIsNoUser(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	testMail := mail.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}
	creator := testee.buildCreator(testMail)

	// Overwrite the mocked interface which are needed in this test:
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", testMail.Sender).Return(false, -1)

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicket", testMail.Subject, creator, mock.Anything).
		Return(&ticket.Ticket{}, nil)

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
	Handling a mail from a user which is registered, should be possible.
*/
func TestIncomingMailHandler_handleIncomingMails_SenderIsAUser(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	userId := 1
	testMail := mail.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}
	testMails := []mail.Mail{testMail}

	// Overwrite the mocked interface which are needed in this test:
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", testMail.Sender).Return(true, userId)
	returnedUser := user.User{}
	mockedUserContext.On("GetUserById", userId).Return(true, returnedUser)

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicketForInternalUser", testMail.Subject, returnedUser, mock.Anything).
		Return(&ticket.Ticket{}, nil)

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
	Handling a mail for when the sender is no user, should be possible.
*/
func TestIncomingMailHandler_handleIncomingMails_SenderIsNoUser(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	testMail := mail.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}
	creator := testee.buildCreator(testMail)
	testMails := []mail.Mail{testMail}

	// Overwrite the mocked interface which are needed in this test:
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", testMail.Sender).Return(false, -1)

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicket", testMail.Subject, creator, mock.Anything).
		Return(&ticket.Ticket{}, nil)

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
	A incoming mail for a existing ticket should notify the creator of the ticket.
*/
func TestIncomingMailHandler_handleIncomingMails_TicketExists_MailSent(t *testing.T) {
	testee := getTestHandlerWithMockedData()
	testMail := mail.Mail{Id: "TestId01", Sender: "testsender@test.de", Receiver: "test@test.de", Subject: "Ticket<1> TestSubject",
		Content: "TestContent", SentTime: time.Now().Unix()}
	testMails := []mail.Mail{testMail}
	existingTicket := ticket.Ticket{}

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("GetTicketById", 1).Return(true, &existingTicket)
	mockedTicketContext.On("AppendMessageToTicket", 1, mock.Anything).
		Return(&existingTicket, nil)

	mockedMailContext := new(mockedForTests.MockedMailContext)
	expectedSubject := "New Entry for your ticket: " + html.EscapeString(testMail.Subject)
	expectedMailContent := mail.BuildNotificationMailContent(existingTicket.Info().Creator.Mail, testMail.Sender, testMail.Content)

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
	mockedTicketContext.On("GetTicketById", mock.Anything).Return(true, &ticket.Ticket{})
	mockedTicketContext.On("AppendMessageToTicket", mock.Anything, mock.Anything).
		Return(&ticket.Ticket{}, errors.New("TestError"))
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
	Assert a MessageEntry ignoring the autonatically set timestamp.
*/
func assertMessageEntryIgnoringTime(t *testing.T, expected ticket.MessageEntry, actual ticket.MessageEntry) {
	assert.Equal(t, expected.Content, actual.Content, "Content should be equal")
	assert.Equal(t, expected.CreatorMail, actual.CreatorMail, "CreatorMail should be equal")
	assert.Equal(t, expected.OnlyInternal, actual.OnlyInternal, "OnlyInternal should be equal")
	assert.Equal(t, expected.CreatorMail, actual.CreatorMail, "CreatorMail should be equal")
}
