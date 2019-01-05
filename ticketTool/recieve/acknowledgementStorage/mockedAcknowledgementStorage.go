package acknowledgementStorage

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"github.com/stretchr/testify/mock"
)

/*
	A mocked acknowledgement storage.
 */
type MockedAcknowledgementStorage struct {
	mock.Mock
}

/*
	A mocked function.
 */
func (m *MockedAcknowledgementStorage)AppendAcknowledgements(acknowledge []mailData.Acknowledgment)error{
	args:=m.Called(acknowledge)
	return args.Error(0)
}

/*
	A mocked function.
 */
func (m *MockedAcknowledgementStorage) DeleteAcknowledges(delete []mailData.Acknowledgment) error {
	args:=m.Called(delete)
	return args.Error(0)
}

func (m *MockedAcknowledgementStorage) ReadAcknowledgements() ([]mailData.Acknowledgment, error) {
	args:=m.Called()
	return args.Get(0).([]mailData.Acknowledgment),args.Error(1)
}

func (m *MockedAcknowledgementStorage) readDataFromFile() error{
	args:=m.Called()
	return args.Error(0)
}

/*
	A mocked function.
 */
func (m *MockedAcknowledgementStorage) writeDataToFile() error {
	args:=m.Called()
	return args.Error(0)
}