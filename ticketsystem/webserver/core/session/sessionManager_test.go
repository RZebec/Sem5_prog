package session

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"strconv"
	"../helpers"
)



func TestLoginSystem_Initialize_InvalidPath_ErrorReturned(t *testing.T) {
	testee := LoginSystem{}

	err := testee.Initialize("")
	assert.Error(t, err, "path to login data storage can not be a empty string")
}

func TestLoginSystem_Initialize_FolderDoesNotExist_FolderCreated(t *testing.T) {
	testee := LoginSystem{}

	// Creating a temp directory and remove it after the test:
	folderPath, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(folderPath)
	if err != nil {
		t.Error(err)
	}
	// Create a path in the temp directory, the following path to a folder will not exist:
	notExistingFolderPath := path.Join(folderPath, "testDirectory")
	err = testee.Initialize(notExistingFolderPath)
	// No error should occur:
	assert.Nil(t, err)
	assert.DirExists(t, notExistingFolderPath)
	assert.FileExists(t, path.Join(notExistingFolderPath, "loginData.json"))
}


func TestLoginSystem_Initialize_DataFileAlreadyExists_DataIsLoaded(t *testing.T) {
	testee := LoginSystem{}

	// Creating a temp directory and remove it after the test:
	folderPath, err := ioutil.TempDir("", "test")
	//defer os.RemoveAll(folderPath)
	if err != nil {
		t.Error(err)
	}
	// Create a path in the temp directory, the following path to a folder will not exist:
	newFolderPath := path.Join(folderPath, "testDirectory")
	helpers.CreateFolderPath(newFolderPath)
	// Write example data to the file
	sampleDataPath := path.Join(newFolderPath, "loginData.json")
	writeTestDataToFile(t, sampleDataPath)

	err = testee.Initialize(newFolderPath)
	// No error should occur:
	assert.Nil(t, err)

	assert.NotNil(t, testee.cachedUserData)
	assert.NotNil(t, testee.cachedUserData.users)

	testData := getTestLoginData()
	t.Log("TestDataSize" + strconv.Itoa(len(testData)))
	assert.Equal(t, len(testee.cachedUserData.users), len(testData))
	assert.ElementsMatch(t, testee.cachedUserData.users, testData)
}

func writeTestDataToFile(t *testing.T, sampleDataPath string) {
	sampleData := []byte(testLoginData)
	err := ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	assert.Nil(t, err)
}

func getTestLoginData()([]storedUserData){
	sampleData := []byte(testLoginData)
	parsedData := new([]storedUserData)
	json.Unmarshal(sampleData, &parsedData)
	return *parsedData
}

func TestLoginSystem_Login(t *testing.T) {
	assert.True(t, false)
}

func TestLoginSystem_Logout(t *testing.T) {
	assert.True(t, false)
}

func TestLoginSystem_RefreshToken(t *testing.T) {
	assert.True(t, false)
}
func TestLoginSystem_Register_NoDataWasStored_UserIsRegistered(t *testing.T) {
	testee := LoginSystem{}

	// Creating a temp directory and remove it after the test:
	folderPath, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(folderPath)
	if err != nil {
		t.Error(err)
	}
	// Create a path in the temp directory, the following path to a folder will not exist:
	tempFolderPath := path.Join(folderPath, "testDirectory")
	err = testee.Initialize(tempFolderPath)
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
