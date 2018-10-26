package session

import (
	"crypto/rand"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"
)

/*
	The session manager provides functions to keep a session alive or to check if the session is valid.
*/
type SessionManager interface {
	// Check if the current session is valid. Also returns the user of the session.
	SessionIsValid(token string) (isValid bool, userId int, userName string, err error)

	// Refresh the token. The new token should be used for all following request.
	RefreshToken(token string) (newToken string, err error)
}

/*
	Represents the stored data. Used for serialization and storage.
*/
type storedUserData struct {
	UserName   string
	UserId     int
	StoredPass string
	StoredSalt string
}

/*
	Interface for the UserManager. Provides functions to register, login or logout at user.
*/
type UserManager interface {
	// Login in a user.
	Login(userName string, password string) (success bool, authToken string, err error)
	// Register a new user.
	Register(userName string, password string) (success bool, err error)
	// Logout a user.
	Logout(authToken string)
}

/*
	A private struct, used to store the current active sessions in memory.
*/
type inMemorySession struct {
	userId           int
	userName         string
	sessionToken     string
	sessionTimestamp time.Time
}

/*
	The LoginSystem contains all parts to handle the access to user and session data.
*/
type LoginSystem struct {
	fileAccessMutex      sync.Mutex
	loginDataFileName    string
	loginDataFilePath    string
	cachedUserDataMutex  sync.RWMutex
	cachedUserData       []storedUserData
	currentSessions      map[string]inMemorySession
	currentSessionsMutex sync.RWMutex
}

/*
	Initializes the LoginSystem. Uses the provided folder path to store the data.
*/
func (s *LoginSystem) Initialize(folderPath string) (err error) {
	s.setDefaultValues()
	// Validate the provided path:
	if folderPath == "" {
		return errors.New("path to login data storage can not be a empty string")
	}
	// We are going to execute file operations => set the lock:
	s.fileAccessMutex.Lock()
	defer s.fileAccessMutex.Unlock()

	// Create the path to the folder (if necessary):
	s.loginDataFilePath = path.Join(folderPath, s.loginDataFileName)
	pathExists, er := helpers.FilePathExists(folderPath)
	if er != nil {
		return er
	}
	if !pathExists {
		er = helpers.CreateFolderPath(folderPath)
		if er != nil {
			return er
		}
	}
	// Create the storage file (if necessary):
	er = helpers.CreateFileIfNotExists(s.loginDataFilePath)
	if er != nil {
		return er
	}
	// Read the data from the file and fill the cache:
	er = s.readFileAndUpdateCache(s.loginDataFilePath)
	if er != nil {
		return er
	}
	return
}

// Refresh the token. The new token should be used for all following request.
func (s *LoginSystem) RefreshToken(token string) (newToken string, err error) {
	valid, userId, userName, err := s.SessionIsValid(token)
	if valid {
		s.currentSessionsMutex.Lock()
		defer s.currentSessionsMutex.Unlock()
		newToken, err := generateUUID()
		if err != nil {
			return "", err
		}
		s.currentSessions[newToken] = inMemorySession{userName: userName, userId: userId, sessionToken: newToken, sessionTimestamp: time.Now()}
		delete(s.currentSessions, token)
		return newToken, nil
	} else {
		return "", errors.New("unknown session")
	}
}

/*
	Register a new user.
*/
func (s *LoginSystem) Register(userName string, password string) (success bool, err error) {
	if userName == "" {
		// TODO: Validator for username
		return false, errors.New("userName not valid")
	}
	if password == "" {
		// TODO: Validator for password
		return false, errors.New("password not valid")
	}
	// Check if user already exists. There can not be multiple users with the same username:
	if s.checkIfUserExistsOnCache(userName) {
		return false, errors.New("user with this name already exists")
	}
	// Register the new user:
	er := s.registerNewUser(userName, password)
	if er != nil {
		return false, errors.New("could not create new user. reason: " + er.Error())
	}
	return true, nil
}

/*
	Logout a user.
*/
func (s *LoginSystem) Logout(authToken string) {
	s.currentSessionsMutex.Lock()
	defer s.currentSessionsMutex.Unlock()
	delete(s.currentSessions, authToken)
}

// Check if the current session is valid. Also returns the user of the session.
func (s *LoginSystem) SessionIsValid(token string) (isValid bool, userId int, userName string, err error) {
	s.currentSessionsMutex.RLock()

	user, ok := s.currentSessions[token]
	if ok {
		s.currentSessionsMutex.RUnlock()
		if time.Now().Sub(user.sessionTimestamp) > time.Duration(10*time.Minute) {
			s.currentSessionsMutex.Lock()
			delete(s.currentSessions, token)
			s.currentSessionsMutex.Unlock()
			return false, -1, "", nil
		}
		return ok, user.userId, user.userName, nil
	} else {
		s.currentSessionsMutex.RUnlock()
	}
	return false, -1, "", nil
}

/*
	Login in a user.
*/
func (s *LoginSystem) Login(userName string, password string) (success bool, authToken string, err error) {
	s.cachedUserDataMutex.RLock()
	defer s.cachedUserDataMutex.RUnlock()

	for _, v := range s.cachedUserData {
		if strings.ToLower(v.UserName) == strings.ToLower(userName) {
			valid := s.checkUserCredentials(v, password)
			if valid {
				var token, err = s.createSessionForUser(v)
				return true, token, err
			}
		}
	}

	return false, "", errors.New("user not found")
}

/*
	Create a session for a user. Returns the token or an error.
*/
func (s *LoginSystem) createSessionForUser(user storedUserData) (authToken string, err error) {
	s.currentSessionsMutex.Lock()
	defer s.currentSessionsMutex.Unlock()
	token, err := generateUUID()
	if err != nil {
		return "", err
	}
	s.currentSessions[token] = inMemorySession{userName: user.UserName, userId: user.UserId, sessionToken: authToken, sessionTimestamp: time.Now()}
	return token, nil
}

/*
	Generate a UUID.
	Source: https://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language
*/
func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, er := rand.Read(b)
	if er != nil {
		return "", er
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

/*
	Check if the provided password for a user is correct.
*/
func (s *LoginSystem) checkUserCredentials(user storedUserData, providedPassword string) (success bool) {
	// TODO: Adjust with password salting and stuff....
	return user.StoredPass == providedPassword
}

/*
	Check if the user already exists (using the cache). Returns true, if the user exists, false if not.
*/
func (s *LoginSystem) checkIfUserExistsOnCache(userName string) bool {
	// We want to read from the cache => get a read lock:
	s.cachedUserDataMutex.RLock()
	defer s.cachedUserDataMutex.RUnlock()
	if s.cachedUserData == nil {
		return false
	} else {
		for _, v := range s.cachedUserData {
			if strings.ToLower(v.UserName) == strings.ToLower(userName) {
				return true
			}
		}
	}
	return false
}

/*
	Register a new user.
*/
func (s *LoginSystem) registerNewUser(userName string, password string) (err error) {
	s.fileAccessMutex.Lock()
	defer s.fileAccessMutex.Unlock()
	existingData, err := s.readJsonDataFromFile(s.loginDataFilePath)
	if err != nil {
		return err
	}
	generatedLoginData := s.generateLoginData(userName, password, len(existingData))
	newData := append(existingData, generatedLoginData)
	err = s.writeJsonDataToFile(s.loginDataFilePath, newData)
	if err != nil {
		return err
	}
	s.cachedUserDataMutex.Lock()
	defer s.cachedUserDataMutex.Unlock()
	s.cachedUserData = append(s.cachedUserData, generatedLoginData)
	return nil
}

/*
	Generates the login data.
*/
func (s *LoginSystem) generateLoginData(userName string, password string, newId int) (userData storedUserData) {
	// TODO: Adjust with password salting and stuff....
	return storedUserData{UserName: userName, StoredPass: password, StoredSalt: "1234", UserId: newId}
}

/*
	Sets the default values.
*/
func (s *LoginSystem) setDefaultValues() {
	s.loginDataFileName = "loginData.json"
	s.currentSessions = make(map[string]inMemorySession)
}

/*
	Read the json data from the file.
*/
func (s *LoginSystem) readJsonDataFromFile(filePath string) (data []storedUserData, err error) {

	fileValue, err := helpers.ReadAllDataFromFile(filePath)
	if err != nil {
		return nil, err
	}
	s.cachedUserDataMutex.Lock()
	defer s.cachedUserDataMutex.Unlock()
	parsedData := new([]storedUserData)
	json.Unmarshal(fileValue, &parsedData)
	return *parsedData, nil
}

/*
	Write the data as json to the file.
*/
func (s *LoginSystem) writeJsonDataToFile(filePath string, data []storedUserData) (err error) {

	jsonData, err := json.MarshalIndent(data, "", "    ")
	return helpers.WriteDataToFile(filePath, jsonData)
}

/*
	Read the data file and update the cached values.
*/
func (s *LoginSystem) readFileAndUpdateCache(filePath string) (err error) {

	parsedData, err := s.readJsonDataFromFile(filePath)
	if parsedData != nil {
		s.cachedUserData = parsedData
	} else {
		s.cachedUserData = *new([]storedUserData)
	}
	return
}
