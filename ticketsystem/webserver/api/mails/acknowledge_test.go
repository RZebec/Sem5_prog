package mails

import (
	"bytes"
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
type MockedAcknowledgementMailContext struct {
	mock.Mock
	ReceivedAcks []mail.Acknowledgment
	throwError   bool
}

/*
	Get test data.
*/
func getTestAcknowledgements() []mail.Acknowledgment {
	var acks []mail.Acknowledgment
	acks = append(acks, mail.Acknowledgment{Id: "id01", Subject: "Subject01"})
	acks = append(acks, mail.Acknowledgment{Id: "id02", Subject: "Subject02"})
	acks = append(acks, mail.Acknowledgment{Id: "id03", Subject: "Subject03"})
	return acks
}

/*
	A mocked function.
*/
func (m *MockedAcknowledgementMailContext) GetUnsentMails() ([]mail.Mail, error) {
	return []mail.Mail{}, nil
}

/*
	A mocked function.
*/
func (m *MockedAcknowledgementMailContext) AcknowledgeMails(acknowledgments []mail.Acknowledgment) error {
	if m.throwError {
		return errors.New("Test error")
	}
	m.ReceivedAcks = acknowledgments
	return nil
}

/*
	A mocked function.
*/
func (m *MockedAcknowledgementMailContext) CreateNewOutgoingMail(receiver string, subject string, content string) error {
	return nil
}

/*
	A valid request should return 200.
*/
func TestAcknowledgeMailHandler_ServeHTTP_OkReturned(t *testing.T) {
	mockedMailContext := new(MockedAcknowledgementMailContext)
	mockedMailContext.throwError = false
	testee := AcknowledgeMailHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext}

	jsonData, _ := json.Marshal(getTestAcknowledgements())

	req, err := http.NewRequest("POST", shared.AcknowledgmentPath, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	for idx, ackData := range getTestAcknowledgements() {
		actual := mockedMailContext.ReceivedAcks[idx]
		assert.Equal(t, ackData, actual, "Acknowledgment should be delivered")
	}
}

/*
	A error should result in a 500.
*/
func TestAcknowledgeMailHandler_ServeHTTP_500Returned(t *testing.T) {
	mockedMailContext := new(MockedAcknowledgementMailContext)
	mockedMailContext.throwError = true
	testee := AcknowledgeMailHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext}

	jsonData, _ := json.Marshal(getTestAcknowledgements())

	req, err := http.NewRequest("POST", shared.AcknowledgmentPath, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")
}

/*
	A invalid payload should result in a 400.
*/
func TestAcknowledgeMailHandler_ServeHTTP_InvalidPayload_400Returned(t *testing.T) {
	mockedMailContext := new(MockedAcknowledgementMailContext)
	mockedMailContext.throwError = false
	testee := AcknowledgeMailHandler{Logger: testhelpers.GetTestLogger(), MailContext: mockedMailContext}

	jsonData, _ := json.Marshal(getTestAcknowledgements())
	// Make the jsonData invalid:
	jsonData[2] = 4

	req, err := http.NewRequest("POST", shared.AcknowledgmentPath, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 400, rr.Code, "Status code 400 should be returned")
}
