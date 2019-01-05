package confirm

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"github.com/stretchr/testify/mock"
)

type MockedConfirm struct {
	mock.Mock
}

func (m *MockedConfirm) GetAllAcknowledges(mails []mailData.Mail) []mailData.Acknowledgment {
	args := m.Called(mails)
	return args.Get(0).([]mailData.Acknowledgment)
}

func (m *MockedConfirm) GetSingleAcknowledges(allAcknowledges []mailData.Acknowledgment, answer string) ([]mailData.Acknowledgment, []mailData.Acknowledgment) {
	args := m.Called(allAcknowledges, answer)
	return args.Get(0).([]mailData.Acknowledgment), args.Get(1).([]mailData.Acknowledgment)
}

func (m *MockedConfirm) ShowAllEmailAcks(allAcknowledges []mailData.Acknowledgment) {
	m.Called(allAcknowledges)
}
