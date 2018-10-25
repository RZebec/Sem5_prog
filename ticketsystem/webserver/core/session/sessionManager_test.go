package session

import (
	"../helpers"
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
)

func prepareTempDirectory() (string, string, error){
	// Creating a temp directory and remove it after the test:
	rootPath, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(rootPath)
	if err != nil {
		return "", "", err
	}
	// Create a path in the temp directory, the following path to a folder will not exist:
	notExistingFolderPath := path.Join(rootPath, "testDirectory")
	return notExistingFolderPath, rootPath,	nil
}


func TestLoginSystem_Initialize_InvalidPath_ErrorReturned(t *testing.T) {
	testee := LoginSystem{}

	err := testee.Initialize("")
	assert.Error(t, err, "path to login data storage can not be a empty string")
}

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

func writeTestDataToFile(t *testing.T, filePath string) {
	os.MkdirAll(filepath.Dir(filePath), 0644)
	sampleData := []byte(testLoginData)
	err := ioutil.WriteFile(filePath, sampleData, 0644)
	assert.Nil(t, err)
}

func getTestLoginData()([]storedUserData){
	sampleData := []byte(testLoginData)
	parsedData := new([]storedUserData)
	json.Unmarshal(sampleData, &parsedData)
	return *parsedData
}

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

func TestLoginSystem_Login_IncrrectLoginData_NotLoggedIn(t *testing.T) {
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

	// The user should no be logged in
	success, token, err := testee.Login(userName, "test123")
	assert.False(t, success, "User should not be logged in")
	assert.Equal(t, "", token, "The token should not be empty")

	// A session should not be created
	sessions := testee.currentSessions
	_, ok := sessions[token]
	assert.False(t, ok)
}

func TestLoginSystem_Logout(t *testing.T) {
	assert.True(t, false)
}

func TestLoginSystem_RefreshToken(t *testing.T) {
	assert.True(t, false)
}
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

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(4)
	// Register some users in a concurrent way:
	go func() {
		defer waitGroup.Done()
		testee.Register("testUser1", "testpassword1")
	} ()
	go func() {
		defer waitGroup.Done()
		testee.Register("testUser2", "testpassword1")
	} ()
	go func() {
		defer waitGroup.Done()
		testee.Register("testUser3", "testpassword1")
	} ()
	go func() {
		defer waitGroup.Done()
		testee.Register("testUser4", "testpassword1")
	} ()

	// wait for all go routines to finish
	waitGroup.Wait()
	writtenData, err := getDataFromFile(testee.loginDataFilePath)
	assert.Nil(t, err)

	assert.Equal(t, len(writtenData), 4)
	for i := 1; i <= 4; i++  {
		found := false
		for _, v := range writtenData {
			if strings.ToLower(v.UserName) == strings.ToLower("testUser" + strconv.Itoa(i)) {
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}
func TestLoginSystem_SessionIsValid(t *testing.T) {
	assert.True(t, false)
}

func getDataFromFile(filePath string) (data []storedUserData, err error){
	fileData, err := helpers.ReadAllDataFromFile(filePath)
	if err != nil {
		return nil, err
	}
	parsedData := new([]storedUserData)
	json.Unmarshal(fileData, &parsedData)
	return *parsedData, nil
}

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
