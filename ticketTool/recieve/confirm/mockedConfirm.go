package confirm

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"github.com/stretchr/testify/mock"
)

type MockedConfirm struct {
	mock.Mock
}

func (m *MockedConfirm) GetAllAcknowledges(mails []mail.Mail) []mail.Acknowledgment {
	args := m.Called(mails)
	return args.Get(0).([]mail.Acknowledgment)
}

func (m *MockedConfirm) GetSingleAcknowledges(allAcknowledges []mail.Acknowledgment, answer string) ([]mail.Acknowledgment, []mail.Acknowledgment) {
	args := m.Called(allAcknowledges, answer)
	return args.Get(0).([]mail.Acknowledgment), args.Get(1).([]mail.Acknowledgment)
}

func (m *MockedConfirm) ShowAllEmailAcks(allAcknowledges []mail.Acknowledgment) {
	m.Called(allAcknowledges)
}
