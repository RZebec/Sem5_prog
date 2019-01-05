package mailGeneration

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"github.com/stretchr/testify/mock"
)

type MockedMailGenerator struct {
	mock.Mock
}

func (m *MockedMailGenerator) RandomMail(n int, subjectLength int, contentLength int) []mailData.Mail {
	args := m.Called(n, subjectLength, contentLength)
	return args.Get(0).([]mailData.Mail )
}
func (m *MockedMailGenerator) ExplicitMail() []mailData.Mail {
	args := m.Called()
	return args.Get(0).([]mailData.Mail )
}