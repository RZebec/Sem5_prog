package client

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"github.com/stretchr/testify/mock"
)

/*
	A mocked client.
*/
type MockedClient struct {
	mock.Mock
}

/*
	A mocked function.
*/
func (m *MockedClient) SendMails(mails []mailData.Mail) error {
	args := m.Called(mails)
	return args.Error(0)
}

/*
	A mocked function.
*/
func (m *MockedClient) ReceiveMails() ([]mailData.Mail, error) {
	args := m.Called()
	return args.Get(0).([]mailData.Mail), args.Error(1)
}

/*
	A mocked function.
*/
func (m *MockedClient) AcknowledgeMails(mailsToAcknowledge []mailData.Acknowledgment) error {
	args := m.Called(mailsToAcknowledge)
	return args.Error(0)
}
