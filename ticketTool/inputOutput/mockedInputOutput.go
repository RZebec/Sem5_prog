// 5894619, 6720876, 9793350
package inputOutput

import "github.com/stretchr/testify/mock"

/*
	A mocked input and output.
*/
type MockedInputOutput struct {
	mock.Mock
}

/*
	A mocked function.
*/
func (m *MockedInputOutput) ReadEntry() string {
	args := m.Called()
	return args.String(0)
}

/*
	A mocked function.
*/
func (m *MockedInputOutput) Print(text string) {
	m.Called(text)
}
