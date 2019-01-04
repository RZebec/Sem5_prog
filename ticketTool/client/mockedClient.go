package client

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"github.com/stretchr/testify/mock"
)

type MockedClient struct {
	mock.Mock
}


func (m *MockedClient) SendMails(mails []mail.Mail) error {
	args := m.Called(mails)
	return args.Error(0)
}
func (m *MockedClient) ReceiveMails() ([]mail.Mail, error) {
	args := m.Called()
	return args.Get(0).([]mail.Mail), args.Error(1)
}
func (m *MockedClient) AcknowledgeMails(mailsToAcknowledge []mail.Acknowledgment) error {
	args := m.Called(mailsToAcknowledge)
	return  args.Error(0)
}