package mockedForTests

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"github.com/stretchr/testify/mock"
)

/*
	A mocked user context.
*/
type MockedUserContext struct {
	mock.Mock
}

/*
	Mocked function.
*/
func (m *MockedUserContext) SessionIsValid(token string) (isValid bool, userId int, userName string, role user.UserRole, err error) {
	args := m.Called(token)
	return args.Bool(0), args.Int(1), args.String(2),
		args.Get(3).(user.UserRole), args.Error(4)
}

/*
	Mocked function.
*/
func (m *MockedUserContext) RefreshToken(token string) (newToken string, err error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedUserContext) Login(userName string, password string) (success bool, authToken string, err error) {
	args := m.Called(userName, password)
	return args.Bool(0), args.String(1), args.Error(2)
}

/*
	Mocked function.
*/
func (m *MockedUserContext) Register(userName string, password string, firstName string, lastName string) (success bool, err error) {
	args := m.Called(userName, password, firstName, lastName)
	return args.Bool(0), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedUserContext) Logout(authToken string) {
	m.Called(authToken)
	return
}

/*
	Mocked function.
*/
func (m *MockedUserContext) EnableVacationMode(token string) (err error) {
	args := m.Called(token)
	return args.Error(0)
}

/*
	Mocked function.
*/
func (m *MockedUserContext) DisableVacationMode(token string) (err error) {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockedUserContext) UnlockAccount(currentUserToken string, userIdToUnlock int) (unlocked bool, err error) {
	args := m.Called(currentUserToken, userIdToUnlock)
	return args.Bool(0), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedUserContext) ChangePassword(currentUserToken string, currentUserPassword string, newPassword string) (changed bool, err error) {
	args := m.Called(currentUserToken, currentUserPassword, newPassword)
	return args.Bool(0), args.Error(1)
}

func (m *MockedUserContext) GetAllLockedUsers() []user.User {
	args := m.Called()
	return args.Get(0).([]user.User)
}

/*
	Mocked function.
*/
func (m *MockedUserContext) GetUserForEmail(mailAddress string) (isRegisteredUser bool, userId int) {
	args := m.Called(mailAddress)
	return args.Bool(0), args.Int(1)
}

/*
	Mocked function.
*/
func (m *MockedUserContext) GetUserById(userId int) (exists bool, existingUser user.User) {
	args := m.Called(userId)
	return args.Bool(0), args.Get(1).(user.User)
}

/*
	Mocked function.
*/
func (m *MockedUserContext) GetAllActiveUsers() []user.User {
	args := m.Called()
	return args.Get(0).([]user.User)
}
