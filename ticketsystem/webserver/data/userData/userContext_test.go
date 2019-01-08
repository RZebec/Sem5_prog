// 5894619, 6720876, 9793350
package userData

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

/*
	Example for the initialization.
*/
func ExampleLoginSystem_Initialize() {
	loginSystem := LoginSystem{}
	loginSystem.Initialize("pathToFolderToUse")
}

/*
	Example for the registration of a user.
*/
func ExampleLoginSystem_Register() {
	loginSystem := LoginSystem{}
	// The path will be created if it does not exist, but this is only an example, so it will be removed again.
	pathToFolderWhichShouldBeUsed := "pathToFolder/subFolder/subSubFolder"
	defer os.RemoveAll(pathToFolderWhichShouldBeUsed)

	loginSystem.Initialize(pathToFolderWhichShouldBeUsed)

	loginSystem.Register("UserName", "UserPassword", "firstName", "lastName")
}

/*
	Example for the login.
*/
func ExampleLoginSystem_Login() {
	loginSystem := LoginSystem{}
	loginSystem.Initialize("pathToFolderToUse")

	loginSystem.Login("UserName", "UserPassword")
}

/*
	Example for the logout.
*/
func ExampleLoginSystem_Logout() {
	loginSystem := LoginSystem{}
	loginSystem.Initialize("pathToFolderToUse")

	token := "1563.....534sf2"

	loginSystem.Logout(token)
}

/*
	Example for the refresh of a token.
*/
func ExampleLoginSystem_RefreshToken() {
	loginSystem := LoginSystem{}
	loginSystem.Initialize("pathToFolderToUse")

	token := "1563.....534sf2"
	// The token will be changed:
	token, _ = loginSystem.RefreshToken(token)
}

/*
	Example to check if a session is valid.
*/
func ExampleLoginSystem_SessionIsValid() {
	loginSystem := LoginSystem{}
	loginSystem.currentSessions = make(map[string]inMemorySession)

	token := "1234567899+op"
	loginSystem.currentSessions[token] = inMemorySession{userId: 1, userMail: "max.mustermann@mail.com",
		userRole: 1, sessionToken: token, sessionTimestamp: time.Now()}

	valid, userId, userName, userRole, _ := loginSystem.SessionIsValid(token)
	fmt.Println(valid)
	fmt.Println(userId)
	fmt.Println(userName)
	fmt.Println(userRole)
	// Output:
	// true
	// 1
	// max.mustermann@mail.com
	// 1

}

/*
	Initializing the login system with an invalid path should return an error.
*/
func TestLoginSystem_Initialize_InvalidPath_ErrorReturned(t *testing.T) {
	testee := LoginSystem{}

	err := testee.Initialize("")
	assert.Error(t, err, "path to login data storage can not be a empty string")
}

/*
	Initializing the login system with a path, which does not exits, should create the folder.
*/
func TestLoginSystem_Initialize_FolderDoesNotExist_FolderCreated(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Error(err)
	}
	err = testee.Initialize(folderPath)
	// No error should occur:
	assert.Nil(t, err)
	assert.DirExists(t, folderPath)
	assert.FileExists(t, path.Join(folderPath, "loginData.json"))
}

/*
	Initializing the login system with an already existing data file, should load the data.
*/
func TestLoginSystem_Initialize_DataFileAlreadyExists_DataIsLoaded(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)
	// No error should occur:
	assert.Nil(t, err)

	assert.NotNil(t, testee.cachedUserData)

	testData := getTestLoginData()
	assert.Equal(t, len(testData), len(testee.cachedUserData))
	assert.ElementsMatch(t, testee.cachedUserData, testData)
}

/*
	A user should be able to change his own password.
*/
func TestLoginSystem_ChangePassword_PasswordChanged(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)
	// No error should occur:
	assert.Nil(t, err)

	userName := "testUser4@test.de"
	userPassword := "asdSsdsdf!1osp"
	newPassword := "Test1234!"
	// The user should be logged in
	success, token, err := testee.Login(userName, userPassword)
	assert.True(t, success, "User should not logged in")
	assert.NotEmpty(t, token, "The token should not be empty")

	changed, err := testee.ChangePassword(token, userPassword, newPassword)
	assert.True(t, changed, "The password should be changed")
	assert.Nil(t, err)
}

/*
	A user should be able to change his own password.
*/
func TestLoginSystem_ChangePassword_InvalidPassword_PasswordNotChanged(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)
	// No error should occur:
	assert.Nil(t, err)

	userName := "testUser4@test.de"
	oldPassword := "asdSsdsdf!1osp"
	newPassword := "Testwew1234!"
	// The user should be logged in
	success, token, err := testee.Login(userName, oldPassword)
	assert.True(t, success, "User should be logged in")
	assert.NotEmpty(t, token, "The token should not be empty")

	// Use a wrong password
	changed, err := testee.ChangePassword(token, oldPassword+"545", newPassword)
	assert.False(t, changed, "The password should not be changed")
	assert.Equal(t, "user password could not be changed", err.Error())
}

/*
	Logging in with a valid login, should result in a new session.
*/
func TestLoginSystem_Login_CorrectLoginData_LoggedIn(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	err = testee.Initialize(folderPath)
	assert.Nil(t, err)

	userName := "testUser@test.de"
	password := "Secret2?"
	firstName := "max"
	lastName := "muster"
	success, err := testee.Register(userName, password, firstName, lastName)

	assert.True(t, success, "User should be registered")
	assert.Nil(t, err)

	success, token, err := testee.Login("Admin@Admin.de", "ChangeMe2018!")
	assert.True(t, success, "Admin should be logged in")

	createdUserId := -1
	for _, entry := range testee.cachedUserData {
		if entry.Mail == userName {
			createdUserId = entry.UserId
			break
		}
	}
	unlocked, err := testee.UnlockAccount(token, createdUserId)
	assert.True(t, unlocked, "user should be unlocked")

	success, token, err = testee.Login(userName, password)

	assert.True(t, success, "User should be logged in")
	assert.NotEmpty(t, token, "The token should not be empty")
	sessions := testee.currentSessions
	_, ok := sessions[token]
	assert.True(t, ok)
}

/*
	Logging in with a valid login, but the account has not been unlocked.
	Should result in a failed login.
*/
func TestLoginSystem_Login_CorrectLoginData_LockedAccount_NotLoggedIn(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	err = testee.Initialize(folderPath)
	assert.Nil(t, err)

	userName := "testUser@test.de"
	password := "seCret25?!"
	firstName := "max"
	lastName := "muster"
	success, err := testee.Register(userName, password, firstName, lastName)

	assert.True(t, success, "User should be registered")
	assert.Nil(t, err)

	success, token, err := testee.Login(userName, password)

	assert.False(t, success, "User should not be logged in")
	assert.Empty(t, token, "The token should be empty")
	sessions := testee.currentSessions
	_, ok := sessions[token]
	assert.False(t, ok)
}

/*
	Login with invalid login data should not be possible.
*/
func TestLoginSystem_Login_IncorrectLoginData_NotLoggedIn(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	err = testee.Initialize(folderPath)
	assert.Nil(t, err)

	userName := "testUser@test.de"
	password := "secrößTet!2"
	firstName := "max"
	lastName := "muster"
	success, err := testee.Register(userName, password, firstName, lastName)

	assert.True(t, success, "User should be registered")
	assert.Nil(t, err)

	// The user should not be logged in
	success, token, err := testee.Login(userName, "test123")
	assert.False(t, success, "User should not be logged in")
	assert.Equal(t, "", token, "The token should be empty")

	// A session should not be created
	sessions := testee.currentSessions
	_, ok := sessions[token]
	assert.False(t, ok)
}

/*
	Logout from a existing session, should logout the user and remove the session.
*/
func TestLoginSystem_Logout_SessionExists_UserLoggedOut(t *testing.T) {
	testee := LoginSystem{}
	testee.currentSessions = make(map[string]inMemorySession)

	token := "1234567899+op"
	testee.currentSessions[token] = inMemorySession{userId: 1, userMail: "testUser@test.de", userRole: RegisteredUser,
		sessionToken: "1234567899+op", sessionTimestamp: time.Now()}

	valid, _, _, _, err := testee.SessionIsValid(token)
	assert.True(t, valid, "User should have a session for this test")
	assert.Nil(t, err)

	testee.Logout(token)

	valid, _, _, _, err = testee.SessionIsValid(token)
	assert.False(t, valid, "User should be logged out")
}

/*
	Refreshing a valid token should be possible.
*/
func TestLoginSystem_RefreshToken_ValidToken_TokenIsRefreshed(t *testing.T) {
	testee := LoginSystem{}
	testee.currentSessions = make(map[string]inMemorySession)

	token := "1234567899+op"
	testee.currentSessions[token] = inMemorySession{userId: 1, userMail: "testUser@test.de",
		sessionToken: token, sessionTimestamp: time.Now()}

	newToken, _ := testee.RefreshToken(token)
	assert.NotEqual(t, token, newToken, "The new token should not be equal to the old token")

	_, contains := testee.currentSessions[token]
	assert.False(t, contains, "The old token should be removed")
}

/*
	Refreshing a invalid token should not be possible.
*/
func TestLoginSystem_RefreshToken_UnknownToken_ErrorReturned(t *testing.T) {
	testee := LoginSystem{}
	testee.currentSessions = make(map[string]inMemorySession)

	token := "1234567899+op"
	testee.currentSessions[token] = inMemorySession{userId: 1, userMail: "testUser@test.de",
		sessionToken: token, sessionTimestamp: time.Now()}

	newToken, err := testee.RefreshToken("1234")
	assert.Error(t, err, "unknown session")
	assert.Equal(t, newToken, "", "token should be empty")
}

/*
	Register a user when no previous data was stored, should be possible.
*/
func TestLoginSystem_Register_NoDataWasStored_UserIsRegistered(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Error(err)
	}
	err = testee.Initialize(folderPath)
	// No error should occur:
	assert.Nil(t, err)

	success, err := testee.Register("testUser1@web.de", "T24=23!sds", "max", "muster")
	assert.True(t, success, "user should be registered")
	assert.Nil(t, err)

	writtenData, err := getDataFromFile(testee.loginDataFilePath)
	assert.Nil(t, err)

	// There should be two accounts. The newly registered and the default admin account.
	assert.Equal(t, 2, len(writtenData))
	assert.Equal(t, "testUser1@web.de", writtenData[1].Mail)
	assert.Equal(t, "max", writtenData[1].FirstName)
	assert.Equal(t, RegisteredUser, writtenData[1].Role)
	assert.Equal(t, WaitingToBeUnlocked, writtenData[1].State)
}

/*
	Register a user when no previous data was stored, should be possible.
*/
func TestLoginSystem_Register_NoDataWasStored_DefaultAccountIsCreated(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Error(err)
	}
	err = testee.Initialize(folderPath)
	// No error should occur:
	assert.Nil(t, err)

	writtenData, err := getDataFromFile(testee.loginDataFilePath)
	assert.Nil(t, err)

	// There should be two accounts. The newly registered and the default admin account.
	assert.Equal(t, 1, len(writtenData))
	assert.Equal(t, "Admin@Admin.de", writtenData[0].Mail)
	assert.Equal(t, "AdminUser", writtenData[0].FirstName)
	assert.Equal(t, "AdminUser", writtenData[0].LastName)
	assert.Equal(t, Admin, writtenData[0].Role)
	assert.Equal(t, Active, writtenData[0].State)
}

/*
	Checking if a session is valid, should return true, when the session is valid.
*/
func TestLoginSystem_SessionIsValid_IsValid(t *testing.T) {
	testee := LoginSystem{}
	testee.currentSessions = make(map[string]inMemorySession)

	token := "1234567899+op"
	testee.currentSessions[token] = inMemorySession{userId: 1, userMail: "testUser@test.de",
		sessionToken: token, sessionTimestamp: time.Now()}

	valid, _, _, _, err := testee.SessionIsValid(token)
	assert.True(t, valid, "User should have a session for this test")
	assert.Nil(t, err)
}

/*
	Checking if a session is valid, should return false, when the session is invalid.
*/
func TestLoginSystem_SessionIsValid_IsInValid(t *testing.T) {
	testee := LoginSystem{}
	testee.currentSessions = make(map[string]inMemorySession)

	token := "1234567899+op"
	testee.currentSessions[token] = inMemorySession{userId: 1, userMail: "testUser@test.de",
		sessionToken: "1234567899+op", sessionTimestamp: time.Now()}

	valid, _, _, _, err := testee.SessionIsValid("5698456")
	assert.False(t, valid, "Session should be invalid")
	assert.Nil(t, err)
}

/*
	A session should be automatically timed out.
*/
func TestLoginSystem_SessionIsValid_SessionTimedOut(t *testing.T) {
	testee := LoginSystem{}
	testee.currentSessions = make(map[string]inMemorySession)

	token := "1234567899+op"
	// Set last update of session to 12 minutes ago.
	timestamp := time.Now().Add(time.Duration(-12) * time.Minute)
	testee.currentSessions[token] = inMemorySession{userId: 1, userMail: "testUser@test.de",
		sessionToken: "1234567899+op", sessionTimestamp: timestamp}

	valid, _, _, _, err := testee.SessionIsValid(token)
	assert.False(t, valid, "Session should be invalid")
	assert.Nil(t, err)
}

/*
	It should be possible to register multiple users in a concurrent way.
*/
func TestLoginSystem_Register_ConcurrentAccess_AllRegistered(t *testing.T) {
	// 250 concurrent user registrations are not expected, but his test is to
	// ensure that the system is able to handle that.
	numberOfRegistrations := 250
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Error(err)
	}
	err = testee.Initialize(folderPath)
	// No error should occur:
	assert.Nil(t, err)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(numberOfRegistrations)
	for i := 0; i < numberOfRegistrations; i++ {

		go func(number int) {
			defer waitGroup.Done()
			testee.Register("testUser"+strconv.Itoa(number)+"@test.de", "!Test23password"+strconv.Itoa(number),
				"firstName"+strconv.Itoa(number), "lastName"+strconv.Itoa(number))
		}(i)
	}

	// wait for all go routines to finish
	waitGroup.Wait()
	writtenData, err := getDataFromFile(testee.loginDataFilePath)
	assert.Nil(t, err)

	// All users should be registered, but there is also the default admin account.
	assert.Equal(t, numberOfRegistrations+1, len(writtenData))
	for i := 0; i < numberOfRegistrations; i++ {
		found := false
		for _, v := range writtenData {
			if strings.ToLower(v.Mail) == strings.ToLower("testUser"+strconv.Itoa(i)+"@test.de") {
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}

/*
	Register a user with an invalid username should not be possible.
*/
func TestLoginSystem_Register_InvalidUsername_NotSuccessful(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Error(err)
	}
	err = testee.Initialize(folderPath)
	assert.Nil(t, err)

	success, err := testee.Register("", "1234", "max", "muster")
	assert.False(t, success, "register operation should not be successful")
	assert.Error(t, err, "userName should be invalid")
	assert.Equal(t, "userName not valid", err.Error())
}

/*
	Register a user with an invalid password should not be possible.
*/
func TestLoginSystem_Register_InvalidPassword_NotSuccessful(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Error(err)
	}
	err = testee.Initialize(folderPath)
	assert.Nil(t, err)

	success, err := testee.Register("testuser@web.de", "", "max", "muster")
	assert.False(t, success, "register operation should not be successful")
	assert.Error(t, err, "password requirements not met")
	assert.Equal(t, "password requirements not met", err.Error())
}

/*
	Register a user with an invalid first name should not be possible.
*/
func TestLoginSystem_Register_InvalidFirstName_NotSuccessful(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Error(err)
	}
	err = testee.Initialize(folderPath)
	assert.Nil(t, err)

	success, err := testee.Register("testuser@web.de", "12T?sdsdaw34", "", "muster")
	assert.False(t, success, "register operation should not be successful")
	assert.Error(t, err, "firstName should be invalid")
	assert.Equal(t, "firstName not valid", err.Error())
}

/*
	Register a user with a invalid last name should not be possible.
*/
func TestLoginSystem_Register_InvalidLastName_NotSuccessful(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Error(err)
	}
	err = testee.Initialize(folderPath)
	assert.Nil(t, err)

	success, err := testee.Register("testuser@web.de", "1T1454?!sd234", "max", "")
	assert.False(t, success, "register operation should not be successful")
	assert.Error(t, err, "lastName should be invalid")
	assert.Equal(t, "lastName not valid", err.Error())
}

/*
	Register multiple users with the same user name should not be possible..
*/
func TestLoginSystem_Register_UserNameAlreadyTaken(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	if err != nil {
		t.Error(err)
	}
	err = testee.Initialize(folderPath)
	assert.Nil(t, err)

	userName := "testuser@web.de"
	password := "1T?1asd234"
	firstName := "Peter"
	lastName := "Meier"

	success, err := testee.Register(userName, password, firstName, lastName)
	assert.True(t, success, "register operation should be successful")
	assert.Nil(t, err)

	success, err = testee.Register(userName, password, firstName, lastName)
	assert.False(t, success, "register operation should not be successful")
	assert.NotNil(t, err)
	assert.Equal(t, "user with this name already exists", err.Error())
}

/*
	It should be possible to unlock a account through the admin.
*/
func TestLoginSystem_UnlockAccount_AccountUnlocked(t *testing.T) {
	testee := LoginSystem{}
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	err = testee.Initialize(folderPath)
	assert.Nil(t, err)

	// Register a user. The user should then be in the state of "WaitingToBeUnlocked".
	userName := "testUser@gmx.de"
	password := "secT?1asdret"
	firstName := "max"
	lastName := "muster"
	success, err := testee.Register(userName, password, firstName, lastName)

	found := false
	createdUserId := -1
	for _, entry := range testee.cachedUserData {
		if entry.Mail == userName {
			assert.Equal(t, WaitingToBeUnlocked, entry.State)
			found = true
			createdUserId = entry.UserId
			break
		}
	}
	assert.True(t, found, "User should be registered and state should be set to waiting to be unlocked.")
	assert.True(t, success, "User should be registered")
	assert.Nil(t, err)

	// Login with the admin and unlock the user:
	success, token, err := testee.Login("Admin@Admin.de", "ChangeMe2018!")
	assert.True(t, success, "Admin should be logged in")

	unlocked, err := testee.UnlockAccount(token, createdUserId)
	assert.Nil(t, err)
	assert.True(t, unlocked, "User account should be unlocked")

	// Validate that the user is unlocked:
	found = false
	for _, entry := range testee.cachedUserData {
		if entry.Mail == userName {
			assert.Equal(t, Active, entry.State, "User account state should now be set to active")
			found = true
			createdUserId = entry.UserId
			break
		}
	}
	assert.True(t, found, "User should be unlocked and the cache should be updated.")
}

/*
	Get all locked users should return all locked users.
*/
func TestLoginSystem_GetAllLockedUsers(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = testee.Initialize(folderPath)

	// Register two accounts to create locked accounts:
	success, err := testee.Register("user20@web.de", "1T23?!sdsd4", "Max", "Meier")
	assert.Nil(t, err)
	assert.True(t, success)
	success, err = testee.Register("user30@web.de", "M989$s1234", "Pia", "Müller")
	assert.Nil(t, err)
	assert.True(t, success)

	lockedUsers := testee.GetAllLockedUsers()
	assert.Equal(t, 2, len(lockedUsers), "There should be two locked users")

	firstUserIsReturned := false
	secondUserIsReturned := false
	for _, user := range lockedUsers {
		if user.Mail == "user20@web.de" {
			firstUserIsReturned = true
		}
		if user.Mail == "user30@web.de" {
			secondUserIsReturned = true
		}
	}
	assert.True(t, firstUserIsReturned, "The first created user should be returned")
	assert.True(t, secondUserIsReturned, "The second created user should be returned")
}

/*
	Get all active users should return all active users.
*/
func TestLoginSystem_GetAllActiveUsers(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	activeUsers := testee.GetAllActiveUsers()

	// Assert that all users are returned:
	assert.Equal(t, 6, len(activeUsers))
}

/**
Trying to unlock a account with a user witch is not a admin should fail.
*/
func TestLoginSystem_UnlockAccount_NoAdminRole(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	// Register a user. The user should then be in the state of "WaitingToBeUnlocked".
	userName := "NewlyCreatedTestUser@test.de"
	password := "Pasd081!asd"
	firstName := "max"
	lastName := "muster"
	success, err := testee.Register(userName, password, firstName, lastName)

	found := false
	createdUserId := -1
	for _, entry := range testee.cachedUserData {
		if entry.Mail == userName {
			assert.Equal(t, WaitingToBeUnlocked, entry.State)
			found = true
			createdUserId = entry.UserId
			break
		}
	}
	assert.True(t, found, "User should be registered and state should be set to waiting to be unlocked.")
	assert.True(t, success, "User should be registered")
	assert.Nil(t, err)

	// Log in with a user which is no admin:
	success, token, err := testee.Login("testUser@test.de", "testPasswort756!")
	assert.True(t, success, "User should be logged in")

	// Unlocking should not be possible:
	unlocked, err := testee.UnlockAccount(token, createdUserId)
	assert.NotNil(t, err)
	assert.False(t, unlocked, "User account should not be unlocked")
	assert.Equal(t, "current session has no permission to unlock accounts", err.Error())

	// Assert that the user is not unlocked:
	found = false
	for _, entry := range testee.cachedUserData {
		if entry.UserId == createdUserId {
			assert.Equal(t, WaitingToBeUnlocked, entry.State)
			found = true
			createdUserId = entry.UserId
			break
		}
	}
	assert.True(t, found, "User state should still be set to waiting to be unlocked.")
}

/*
	Unlocking a account witch is not in the correct state should return an error.
*/
func TestLoginSystem_UnlockAccount_AccountInWrongState(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	// Assert that the target user is already unlocked and set to active:
	found := false
	createdUserId := -1
	for _, entry := range testee.cachedUserData {
		if entry.Mail == "testUser3@test.de" {
			assert.Equal(t, Active, entry.State)
			found = true
			createdUserId = entry.UserId
			break
		}
	}
	assert.True(t, found, "User should be registered and state should be set to active.")
	assert.Nil(t, err)

	// Login with the admin to execute the operation:
	success, token, err := testee.Login("testAdmin@test.de", "Hjssdfi=2!ß9!")
	assert.True(t, success, "Admin should be logged in")

	// Assert that the account can not be unlocked again:
	unlocked, err := testee.UnlockAccount(token, createdUserId)
	assert.NotNil(t, err)
	assert.False(t, unlocked, "User account should not be unlocked")
	assert.Equal(t, "can not unlock a account, which is not in the waiting to be unlocked state", err.Error())
}

/*
	Unlocking a account with an unknown id should return a error message.
*/
func TestLoginSystem_UnlockAccount_UnknownAccount(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	// Login with the admin to execute the operation:
	success, token, err := testee.Login("testAdmin@test.de", "Hjssdfi=2!ß9!")
	assert.True(t, success, "Admin should be logged in")

	// Assert that a error message is returned.
	unlocked, err := testee.UnlockAccount(token, 9999)
	assert.NotNil(t, err)
	assert.False(t, unlocked, "User account should not be unlocked")
	assert.Equal(t, "user to unlock not found", err.Error())
}

/*
	Enable the vacation mode for the own account.
*/
func TestLoginSystem_EnableVacationMode_Enabled(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	// Login and assert that the state is set to active:
	success, token, err := testee.Login("testUser5@test.de", "Tzqweq23Aws!")
	assert.True(t, success, "User should be logged in")

	found := false
	for _, entry := range testee.cachedUserData {
		if entry.Mail == "testUser5@test.de" {
			assert.Equal(t, Active, entry.State)
			found = true
			break
		}
	}
	assert.True(t, found, "User should be registered and state should be active.")
	assert.Nil(t, err)

	// Enable the vacation mode:
	err = testee.EnableVacationMode(token)
	assert.Nil(t, err)

	// Assert that the vacation mode has been set.
	found = false
	for _, entry := range testee.cachedUserData {
		if entry.Mail == "testUser5@test.de" {
			assert.Equal(t, OnVacation, entry.State)
			found = true
			break
		}
	}
	assert.True(t, found, "User state should be set to on vacation.")
	assert.Nil(t, err)
}

/*
	Disabling the vacation mode should be possible.
*/
func TestLoginSystem_DisableVacationMode_Disabled(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	// Login and set the vacation mode:
	success, token, err := testee.Login("testUser5@test.de", "Tzqweq23Aws!")
	assert.True(t, success, "User should be logged in")

	found := false
	for _, entry := range testee.cachedUserData {
		if entry.Mail == "testUser5@test.de" {
			assert.Equal(t, Active, entry.State)
			found = true
			break
		}
	}
	assert.True(t, found, "User should be registered and state should be active.")
	assert.Nil(t, err)

	err = testee.EnableVacationMode(token)
	assert.Nil(t, err)

	// Assert that the validation mode is set:
	found = false
	for _, entry := range testee.cachedUserData {
		if entry.Mail == "testUser5@test.de" {
			assert.Equal(t, OnVacation, entry.State)
			found = true
			break
		}
	}
	assert.True(t, found, "User state should be set to on vacation.")
	assert.Nil(t, err)

	// Disable the vacation mode and assert that is has been disabled:
	err = testee.DisableVacationMode(token)
	assert.Nil(t, err)

	found = false
	for _, entry := range testee.cachedUserData {
		if entry.Mail == "testUser5@test.de" {
			assert.Equal(t, Active, entry.State)
			found = true
			break
		}
	}
	assert.True(t, found, "User state should be set to active.")
	assert.Nil(t, err)
}

/*
	Disabling the vacation mode while the user is not in the vacation mode, should return a error message.
*/
func TestLoginSystem_DisableVacationMode_WrongState(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	// Login and validate that the user is not in vacation mode:
	success, token, err := testee.Login("testUser5@test.de", "Tzqweq23Aws!")
	assert.True(t, success, "User should be logged in")

	found := false
	for _, entry := range testee.cachedUserData {
		if entry.Mail == "testUser5@test.de" {
			assert.Equal(t, Active, entry.State)
			found = true
			break
		}
	}
	assert.True(t, found, "User should be registered and state should be active.")
	assert.Nil(t, err)

	// Disable the vacation mode and validate the error message:
	err = testee.DisableVacationMode(token)
	assert.Equal(t, "can not set account to active, when it is not on vacation mode", err.Error())
}

/*
	Getting a existing user should be possible.
*/
func TestLoginSystem_GetUserById(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	exists, user := testee.GetUserById(1)
	assert.True(t, exists, "User should exist")
	assert.Equal(t, "testUser5@test.de", user.Mail, "Correct user should be returned")
}

/*
	Getting a non existing user should return false.
*/
func TestLoginSystem_GetUserById_UserIdDoesNotExist(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	exists, user := testee.GetUserById(9999)
	assert.False(t, exists, "User should not exist")
	assert.Equal(t, "", user.Mail, "Correct user should be returned")
}

/*
	Getting a existing user should be possible.
*/
func TestLoginSystem_GetUserForEmail(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	exists, userId := testee.GetUserForEmail("testUser5@test.de")
	assert.True(t, exists, "User should exist")
	assert.Equal(t, 1, userId, "Correct user id should be returned")
}

/*
	Getting a non existing user should return false.
*/
func TestLoginSystem_GetUserForEmail_UserIdDoesNotExist(t *testing.T) {
	testee := LoginSystem{}

	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)
	// Write example data to the file
	sampleDataPath := path.Join(folderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(folderPath)

	exists, userId := testee.GetUserForEmail("test1254799@ere.de")
	assert.False(t, exists, "User should not exist")
	assert.Equal(t, -1, userId, "Correct user id should be returned")
}

/*
	Preparing a temporary directory for tests, which need access to the file system.
*/
func prepareTempDirectory() (string, string, error) {
	// Creating a temp directory and remove it after the test:
	rootPath, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(rootPath)
	if err != nil {
		return "", "", err
	}
	// Create a path in the temp directory, the following path to a folder will not exist:
	notExistingFolderPath := path.Join(rootPath, "testDirectory")
	return notExistingFolderPath, rootPath, nil
}

/*
	Write the test data to a file.
*/
func writeTestDataToFile(t *testing.T, filePath string) {
	os.MkdirAll(filepath.Dir(filePath), 0644)
	sampleData := []byte(testLoginData)
	err := ioutil.WriteFile(filePath, sampleData, 0644)
	assert.Nil(t, err)
}

/*
	Get the test data.
*/
func getTestLoginData() []storedUserData {
	sampleData := []byte(testLoginData)
	parsedData := new([]storedUserData)
	json.Unmarshal(sampleData, &parsedData)
	return *parsedData
}

/*
	Read the data from a file.
*/
func getDataFromFile(filePath string) (data []storedUserData, err error) {
	fileData, err := helpers.ReadAllDataFromFile(filePath)
	if err != nil {
		return nil, err
	}
	parsedData := new([]storedUserData)
	json.Unmarshal(fileData, &parsedData)
	return *parsedData, nil
}

/*
	The test data to use.
*/
const testLoginData = `[
	{
		"Mail": "testUser@test.de",
		"UserId": 0,
		"FirstName": "Max0",
		"LastName": "Maximum0",
		"StoredPass": "0WPNL5ywpGNBUjiko/55etIliY7NpkSwPf2kk90bcGfwcCHa+KZx7WQNy/C/Z0c2+J/AjBhue1iHt3M4EqXwow==",
		"StoredSalt": "XiiaAjINcE8u/PlMpKffAcuRHD1YgnsbEsrvaO3eo1t6dsSp8eYn7NoH43q6Y96t4a32ZKvjhgp2y9n3amAnAPTmC4DOLISiVOnqDf+/0r0YPl8bauT+cCMsc/bJxvE2",
		"Role": 2,
        "State": 1
	},
	{
		"Mail": "testUser5@test.de",
		"UserId": 1,
		"FirstName": "Max1",
		"LastName": "Maximum1",
		"StoredPass": "NRGv3RuHkLPl0CsVedWwZZe60DY0eMA3UnoOgw879hTBKpl6w6Lz8TgfJQEHSrqjE/J3vHo6Sm50+OhkxrMn3A==",
		"StoredSalt": "ERVP+8Dga01sMRZdqHbkwXza1eBRnSJDGejGMIROd6+z/G8idZpE5z3lxis9f5ZnYYC+4EA8W7+PNa2h23saeQQ2BDkLYZ/3315OJ/NCJk8vGcisrWaO6on19isQGqAp",
		"Role": 2,
        "State": 1
	},
	{
		"Mail": "testUser2@test.de",
		"UserId": 2,
		"FirstName": "Max2",
		"LastName": "Maximum2",
		"StoredPass": "nfSyvQW+W7jS3LcvSOfe9mTlfetFGHNlSh0YDDLPCCQRfVJfyXzWfb6j91FRxbEiIsnVsxWEYp8uSQsrTx83SA==",
		"StoredSalt": "RkG5JQ9leBNzLPhY2B48Ta9mqODLJp1wXkmabcNqjfSra1ITvP+RcFV7SPIwyE8cwIRx2Tbpbnmy9P2kFWOYPiV9nQ3yf6NdiqpnM1FGrysTlqM9Xc9DW6shN23b9GMN",
		"Role": 2,
        "State": 1
	},
	{
		"Mail": "testUser3@test.de",
		"UserId": 3,
		"FirstName": "Max3",
		"LastName": "Maximum3",
		"StoredPass": "NhRxU2BuNYjBbVdPZI/iz0WPtzSZoRgfK06Ekj38fk+UETaMB8Y7GJ0PqlrdFOC8aql25JPV18obGn4d+4H7ag==",
		"StoredSalt": "9mEoaUEQ15isv3GDqHVOZXcwFlNAKNpyN8iT7IrY55fz+7WOlLIQiIwUnEQTGL7pNYdCaIGyBGpFt8wUw6m2FWG+HCM6XmG+CFANG2Gd67Oa0knymEdyz6Mb5l0hVJqw",
		"Role": 2,
        "State": 1
	},
	{
		"Mail": "testUser4@test.de",
		"UserId": 4,
		"FirstName": "Max4",
		"LastName": "Maximum4",
		"StoredPass": "zAgjXlxXpT0MSINlKF/ZsJpT8jkLUtdMc2UG5YU8x9OMxSOk0AfvK0l9AsKX8ErU+KC69AZIML5D76d1B1Oogg==",
		"StoredSalt": "2hMoNTCjIgVaiuE2iko0GijAqk+bGdpthHWNVnNj1PeaCNNl/h6CZ1RuaYQR/M6oWpgqPPWYUP5A3sdnmR9Bc1YGV+1kocOc+NFnKjTwLMnvf+4I2p9031CQwvwChrNY",
		"Role": 2,
        "State": 1
	},
	{
		"Mail": "testAdmin@test.de",
		"UserId": 5,
		"FirstName": "Max5",
		"LastName": "Maximum5",
		"StoredPass": "G+MMAOz7KBZWJCTBW1AJxGWAuiSENOSNx2UVTImdAmdbu2vzEBCuF8qyMppXWv71RxuofdG4RYxULCa+zISBRw==",
		"StoredSalt": "PFeIEH4fO8JLUYwoweahb+LJ0kxgTq4MC3A8nrHx4zCMGFipFqVdf/Uv2yf9ssM+eAeSSOHmV42ngsvhROXCPRHO5kltc0C0M/A5mpAStEqX4IcvvE/1SvblUxnCNbuv",
		"Role": 1,
        "State": 1
	}
]`
