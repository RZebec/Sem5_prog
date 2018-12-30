package mockedForTests

import "github.com/stretchr/testify/mock"

/*
	A mocked api configuration.
 */
type MockedApiConfiguration struct {
	mock.Mock
}

/*
	Mocked function.
*/
func (m *MockedApiConfiguration) GetIncomingMailApiKey() string {
	args := m.Called()
	return args.String(0)
}

/*
	Mocked function.
*/
func (m *MockedApiConfiguration) GetOutgoingMailApiKey() string {
	args := m.Called()
	return args.String(0)
}

/*
	Mocked function.
*/
func (m *MockedApiConfiguration) ChangeIncomingMailApiKey(newKey string) error {
	args := m.Called()
	return args.Error(0)
}

/*
	Mocked function.
*/
func (m *MockedApiConfiguration) ChangeOutgoingMailApiKey(newKey string) error {
	args := m.Called()
	return args.Error(0)
}

/*
	Mocked function.
*/
func (m *MockedApiConfiguration) Validate() (bool, string) {
	args := m.Called()
	return args.Bool(0), args.String(1)
}
