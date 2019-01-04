package acknowledgementStorage

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"github.com/stretchr/testify/mock"
)

type MockedAcknowledementStorage struct {
	mock.Mock
}

func (m *MockedAcknowledementStorage)AppendAcknowledgements(acknowledge []mail.Acknowledgment)error{
	args:=m.Called(acknowledge)
	return args.Error(0)
}

func (m *MockedAcknowledementStorage) DeleteAcknowledges(delete []mail.Acknowledgment) error {
	args:=m.Called(delete)
	return args.Error(0)
}

func (m *MockedAcknowledementStorage) ReadAcknowledgements() ([]mail.Acknowledgment, error) {
	args:=m.Called()
	return args.Get(0).([]mail.Acknowledgment),args.Error(1)
}

func (m *MockedAcknowledementStorage) readDataFromFile() error{
	args:=m.Called()
	return args.Error(0)
}

func (m *MockedAcknowledementStorage) writeDataToFile() error {
	args:=m.Called()
	return args.Error(0)
}