package session

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"encoding/json"
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
	assert.NotNil(t, testee.cachedUserData)

	testData := getTestLoginData()
	t.Log("TestDataSize" + strconv.Itoa(len(testData)))
	assert.Equal(t, len(testData), len(testee.cachedUserData))
	assert.ElementsMatch(t, testee.cachedUserData, testData)
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

	userName := "testUser"
	password := "secret"
	success, err := testee.Register(userName, password)

	assert.True(t, success, "User should be registered")
	assert.Nil(t, err)

	success, token, err := testee.Login(userName, password)
	t.Log(success)
	assert.True(t, success, "User should be logged in")
	assert.NotEmpty(t, token, "The token should not be empty")
	sessions := testee.currentSessions
	_, ok := sessions[token]
	assert.True(t, ok)
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

	userName := "testUser"
	password := "secret"
	success, err := testee.Register(userName, password)

	assert.True(t, success, "User should be registered")
	assert.Nil(t, err)

	// The user should not be logged in
	success, token, err := testee.Login(userName, "test123")
	assert.False(t, success, "User should not be logged in")
	assert.Equal(t, "", token, "The token should not be empty")

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
	testee.currentSessions[token] = inMemorySession{userId: 1, userName: "testUser",
		sessionToken: "1234567899+op", sessionTimestamp: time.Now()}

	valid, _, _, err := testee.SessionIsValid(token)
	assert.True(t, valid, "User should have a session for this test")
	assert.Nil(t, err)

	testee.Logout(token)

	valid, _, _, err = testee.SessionIsValid(token)
	assert.False(t, valid, "User should be logged out")
}

/*
	Refreshing a valid token should be possible.
*/
func TestLoginSystem_RefreshToken_ValidToken_TokenIsRefreshed(t *testing.T) {
	testee := LoginSystem{}
	testee.currentSessions = make(map[string]inMemorySession)

	token := "1234567899+op"
	testee.currentSessions[token] = inMemorySession{userId: 1, userName: "testUser",
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
	testee.currentSessions[token] = inMemorySession{userId: 1, userName: "testUser",
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

	success, err := testee.Register("testUser1", "testpassword1")
	assert.True(t, success, "user should be registered")
	assert.Nil(t, err)

	writtenData, err := getDataFromFile(testee.loginDataFilePath)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(writtenData))
	assert.Equal(t, writtenData[0].UserName, "testUser1")
}

/*
	Checking if a session is valid, should return true, when the session is valid.
*/
func TestLoginSystem_SessionIsValid_IsValid(t *testing.T) {
	testee := LoginSystem{}
	testee.currentSessions = make(map[string]inMemorySession)

	token := "1234567899+op"
	testee.currentSessions[token] = inMemorySession{userId: 1, userName: "testUser",
		sessionToken: token, sessionTimestamp: time.Now()}

	valid, _, _, err := testee.SessionIsValid(token)
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
	testee.currentSessions[token] = inMemorySession{userId: 1, userName: "testUser",
		sessionToken: "1234567899+op", sessionTimestamp: time.Now()}

	valid, _, _, err := testee.SessionIsValid("5698456")
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
	testee.currentSessions[token] = inMemorySession{userId: 1, userName: "testUser",
		sessionToken: "1234567899+op", sessionTimestamp: timestamp}

	valid, _, _, err := testee.SessionIsValid(token)
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
			testee.Register("testUser"+strconv.Itoa(number), "testpassword"+strconv.Itoa(number))
		}(i)
	}

	// wait for all go routines to finish
	waitGroup.Wait()
	writtenData, err := getDataFromFile(testee.loginDataFilePath)
	assert.Nil(t, err)

	assert.Equal(t, numberOfRegistrations, len(writtenData))
	for i := 0; i < numberOfRegistrations; i++ {
		found := false
		for _, v := range writtenData {
			if strings.ToLower(v.UserName) == strings.ToLower("testUser"+strconv.Itoa(i)) {
				found = true
				break
			}
		}
		assert.True(t, found)
	}
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

	success, err := testee.Register("", "1234")
	assert.False(t, success, "register operation should not be successful")
	assert.Error(t, err, "userName should be invalid")
	assert.Equal(t, "userName not valid", err.Error())
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
		"UserName": "testUser",
		"UserId": 0,
		"StoredPass": "testPassword",
		"StoredSalt": "1234"
	},
	{
		"UserName": "testUser",
		"UserId": 1,
		"StoredPass": "testPassword",
		"StoredSalt": "1234"
	},
	{
		"UserName": "testUser3",
		"UserId": 2,
		"StoredPass": "testPassword2",
		"StoredSalt": "1234"
	},
	{
		"UserName": "testUser2",
		"UserId": 3,
		"StoredPass": "testPassword2",
		"StoredSalt": "1234"
	},
	{
		"UserName": "testUser4",
		"UserId": 4,
		"StoredPass": "testPassword2",
		"StoredSalt": "1234"
	},
	{
		"UserName": "testUser5",
		"UserId": 5,
		"StoredPass": "testPassword2",
		"StoredSalt": "1234"
	}
	]`
