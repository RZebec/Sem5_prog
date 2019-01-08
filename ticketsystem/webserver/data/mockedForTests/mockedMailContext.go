// 5894619, 6720876, 9793350
package mockedForTests

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"github.com/stretchr/testify/mock"
)

/*
	A mocked mail context.
*/
type MockedMailContext struct {
	mock.Mock
}

/*
	Mocked function.
*/
func (m *MockedMailContext) GetUnsentMails() ([]mailData.Mail, error) {
	args := m.Called()
	return args.Get(0).([]mailData.Mail), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedMailContext) AcknowledgeMails(acknowledgments []mailData.Acknowledgment) error {
	args := m.Called(acknowledgments)
	return args.Error(0)
}

/*
	Mocked function.
*/
func (m *MockedMailContext) CreateNewOutgoingMail(receiver string, subject string, content string) error {
	args := m.Called(receiver, subject, content)
	return args.Error(0)
}
