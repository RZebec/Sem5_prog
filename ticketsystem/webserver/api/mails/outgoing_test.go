package mails

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	A mocked mail context.
*/
type MockedOutgoingMailContext struct {
	mock.Mock
	ReceivedAcks []mail.Acknowledgment
	throwError   bool
}

/*
	Get test mails.
*/
func getTestMails() []mail.Mail {
	var testMails []mail.Mail
	testMails = append(testMails, mail.Mail{Id: "testId1", Sender: "test@test.de", Receiver: "testReceiver1@test.de",
		Subject: "Ticket<1> testSubject1", Content: "testContent1"})
	testMails = append(testMails, mail.Mail{Id: "testId2", Sender: "test@test.de", Receiver: "testReceiver2@test.de",
		Subject: "testSubject2", Content: "testContent2"})
	testMails = append(testMails, mail.Mail{Id: "testId3", Sender: "test@test.de", Receiver: "testReceiver3@test.de",
		Subject: "testSubject3", Content: "testContent3"})
	return testMails
}

/*
	A mocked function.
*/
func (m *MockedOutgoingMailContext) GetUnsentMails() ([]mail.Mail, error) {
	if m.throwError {
		return []mail.Mail{}, errors.New("Test error")
	}
	return getTestMails(), nil
}

/*
	A mocked function.
*/
func (m *MockedOutgoingMailContext) AcknowledgeMails(acknowledgments []mail.Acknowledgment) error {
	return nil
}

/*
	A mocked function.
*/
func (m *MockedOutgoingMailContext) CreateNewOutgoingMail(receiver string, subject string, content string) error {
	return nil
}

/*
	The provided mails should be in the response of the request.
*/
func TestOutgoingMailHandler_ServeHTTP_MailsReceived(t *testing.T) {
	mockedMailContext := new(MockedOutgoingMailContext)
	mockedMailContext.throwError = false
	testee := OutgoingMailHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext}

	req, err := http.NewRequest("GET", shared.ReceivePath, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	decoder := json.NewDecoder(rr.Body)
	var data []mail.Mail
	err = decoder.Decode(&data)
	assert.Nil(t, err)

	for idx, expectedMail := range getTestMails() {
		actualMail := data[idx]
		assert.Equal(t, expectedMail, actualMail, "The received mails should be equal to the send ones")
	}
}

/*
	A error should return a 500.
*/
func TestOutgoingMailHandler_ServeHTTP_Error_500Returned(t *testing.T) {
	mockedMailContext := new(MockedOutgoingMailContext)
	mockedMailContext.throwError = true
	testee := OutgoingMailHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext}

	req, err := http.NewRequest("GET", shared.ReceivePath, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")
}
