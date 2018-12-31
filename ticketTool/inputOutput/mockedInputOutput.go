package inputOutput

import "github.com/stretchr/testify/mock"

type MockedInputOutput struct {
	mock.Mock
}


func (m *MockedInputOutput) ReadEntry() string {
	args := m.Called()
	return args.String(0)
}
func (m *MockedInputOutput) Print(text string) {
	m.Called(text)
}