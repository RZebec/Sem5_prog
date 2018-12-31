package mailGeneration

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"github.com/stretchr/testify/mock"
)

type MockedMailGenerator struct {
	mock.Mock
}

func (m *MockedMailGenerator) RandomMail(n int, subjectLength int, contentLength int) []mail.Mail {
	args := m.Called(n, subjectLength, contentLength)
	return args.Get(0).([]mail.Mail )
}
func (m *MockedMailGenerator) ExplicitMail() []mail.Mail {
	args := m.Called()
	return args.Get(0).([]mail.Mail )
}