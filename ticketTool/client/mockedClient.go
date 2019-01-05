package client

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"github.com/stretchr/testify/mock"
)

type MockedClient struct {
	mock.Mock
}


func (m *MockedClient) SendMails(mails []mailData.Mail) error {
	args := m.Called(mails)
	return args.Error(0)
}
func (m *MockedClient) ReceiveMails() ([]mailData.Mail, error) {
	args := m.Called()
	return args.Get(0).([]mailData.Mail), args.Error(1)
}
func (m *MockedClient) AcknowledgeMails(mailsToAcknowledge []mailData.Acknowledgment) error {
	args := m.Called(mailsToAcknowledge)
	return  args.Error(0)
}