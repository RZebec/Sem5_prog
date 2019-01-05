/*
	The session package handles the users and the session.
*/
package userData

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/validation/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/validation/passwordRequirements"
	"encoding/json"
	"errors"
	"io"
	"path"
	"strings"
	"sync"
	"time"
)

/*
	The user context provides functions to keep a session alive or to check if the session is valid and to register,
	login or logout at user.
*/
type UserContext interface {
	// Check if the current session is valid. Also returns the user of the session.
	SessionIsValid(token string) (isValid bool, userId int, userName string, role UserRole, err error)
	// Refresh the token. The new token should be used for all following request.
	RefreshToken(token string) (newToken string, err error)
	// Login a user.
	Login(userName string, password string) (success bool, authToken string, err error)
	// Register a new user.
	Register(userName string, password string, firstName string, lastName string) (success bool, err error)
	// Logout a user.
	Logout(authToken string)
	// Set a account to vacation mode. Only possible for the currently logged-in account.
	EnableVacationMode(token string) (err error)
	// Disable the vacation mode. Only possible for the currently logged-in account.
	DisableVacationMode(token string) (err error)
	// Unlock a account which is waiting to be unlocked. The current session needs the permission to do this.
	UnlockAccount(currentUserToken string, userIdToUnlock int) (unlocked bool, err error)
	// Changing the password of a user should be possible, but only for the user himself.
	ChangePassword(currentUserToken string, currentUserPassword string, newPassword string) (changed bool, err error)
	// Get all locked users:
	GetAllLockedUsers() []User
	// Get all active users:
	GetAllActiveUsers() []User
	// Check if the given mail is for a registered user.
	GetUserForEmail(mailAddress string) (isRegisteredUser bool, userId int)
	// Get user by its id.
	GetUserById(userId int) (exists bool, user User)
}

/*
	Represents the stored data. Used for serialization and storage.
*/
type storedUserData struct {
	Mail       string
	UserId     int
	FirstName  string
	LastName   string
	StoredPass []byte
	StoredSalt []byte
	Role       UserRole
	State      UserState
}

/*
	A private struct, used to store the current active sessions in memory.
*/
type inMemorySession struct {
	userId           int
	userMail         string
	userRole         UserRole
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
	var er = s.initializeFiles(folderPath)
	if er != nil {
		return er
	}
	if len(s.cachedUserData) == 0 {
		er := s.registerNewUser("Admin@Admin.de", "ChangeMe2018!",
			"AdminUser", "AdminUser", Admin, Active)
		if er != nil {
			return er
		}
	}

	return
}

/*
	Get a user by id.
*/
func (s *LoginSystem) GetUserById(userId int) (exists bool, user User) {
	s.cachedUserDataMutex.RLock()
	defer s.cachedUserDataMutex.RUnlock()

	for _, userData := range s.cachedUserData {
		if userData.UserId == userId {
			copiedUser := User{Mail: userData.Mail, UserId: userData.UserId, FirstName: userData.FirstName,
				LastName: userData.LastName, Role: userData.Role, State: userData.State}
			copiedUser = copiedUser.Copy()
			return true, copiedUser
		}
	}
	return false, User{}
}

/*
	Get a user from its mail.
*/
func (s *LoginSystem) GetUserForEmail(mailAddress string) (isRegisteredUser bool, userId int) {
	s.cachedUserDataMutex.RLock()
	defer s.cachedUserDataMutex.RUnlock()

	for _, userData := range s.cachedUserData {
		if strings.ToLower(userData.Mail) == strings.ToLower(mailAddress) {
			return true, userData.UserId
		}
	}
	return false, -1
}

// Refresh the token. The new token should be used for all following request.
func (s *LoginSystem) RefreshToken(token string) (newToken string, err error) {
	valid, userId, userName, userRole, err := s.SessionIsValid(token)
	if valid {
		s.currentSessionsMutex.Lock()
		defer s.currentSessionsMutex.Unlock()
		newToken, err := helpers.GenerateUUID()
		if err != nil {
			return "", err
		}
		s.currentSessions[newToken] = inMemorySession{userMail: userName, userId: userId, userRole: userRole,
			sessionToken: newToken, sessionTimestamp: time.Now()}
		delete(s.currentSessions, token)
		return newToken, nil
	} else {
		return "", errors.New("unknown session")
	}
}

/*
	Register a new user.
*/
func (s *LoginSystem) Register(userName string, password string, firstName string, lastName string) (success bool, err error) {
	mailValidator := mail.NewValidator()
	passwordValidator := passwordRequirements.NewValidator()
	if !mailValidator.Validate(userName) {
		return false, errors.New("userName not valid")
	}
	if !passwordValidator.Validate(password) {
		return false, errors.New("password requirements not met")
	}
	if firstName == "" {
		return false, errors.New("firstName not valid")
	}
	if lastName == "" {
		return false, errors.New("lastName not valid")
	}
	// Check if user already exists. There can not be multiple users with the same username:
	if s.checkIfUserExistsOnCache(userName) {
		return false, errors.New("user with this name already exists")
	}
	// Register the new user:
	er := s.registerNewUser(userName, password, firstName, lastName, RegisteredUser, WaitingToBeUnlocked)
	if er != nil {
		return false, errors.New("could not create new user. reason: " + er.Error())
	}
	return true, nil
}

/*
	Change the password for the user of the given session.
*/
func (s *LoginSystem) ChangePassword(currentUserToken string, currentUserPassword string, newPassword string) (changed bool, err error) {
	isValid, userId, userName, _, err := s.SessionIsValid(currentUserToken)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, errors.New("invalid session token")
	}

	for _, v := range s.cachedUserData {
		if strings.ToLower(v.Mail) == strings.ToLower(userName) && v.UserId == userId {
			valid := s.checkUserCredentials(v, currentUserPassword)
			if valid {
				user := User{Mail: v.Mail, UserId: v.UserId, FirstName: v.FirstName, LastName: v.LastName, Role: v.Role,
					State: v.State}
				return s.changeUserPassword(user, newPassword)
			}
		}
	}
	return false, errors.New("user password could not be changed")
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
func (s *LoginSystem) SessionIsValid(token string) (isValid bool, userId int, userName string, userRole UserRole, err error) {
	s.currentSessionsMutex.RLock()

	user, ok := s.currentSessions[token]
	if ok {
		s.currentSessionsMutex.RUnlock()
		if time.Now().Sub(user.sessionTimestamp) > time.Duration(10*time.Minute) {
			s.currentSessionsMutex.Lock()
			delete(s.currentSessions, token)
			s.currentSessionsMutex.Unlock()
			return false, -1, "", -1, nil
		}
		return ok, user.userId, user.userMail, user.userRole, nil
	} else {
		s.currentSessionsMutex.RUnlock()
	}
	return false, -1, "", -1, nil
}

/*
	Login in a user.
*/
func (s *LoginSystem) Login(userName string, password string) (success bool, authToken string, err error) {
	s.cachedUserDataMutex.RLock()
	defer s.cachedUserDataMutex.RUnlock()

	for _, v := range s.cachedUserData {
		if strings.ToLower(v.Mail) == strings.ToLower(userName) {
			if v.State == WaitingToBeUnlocked {
				return false, "", errors.New("user is still waiting to be unlocked")
			}
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
	Setting your own account to vacation mode.
*/
func (s *LoginSystem) DisableVacationMode(token string) (err error) {
	valid, userId, _, _, err := s.SessionIsValid(token)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("can not set vacation mode for invalid session")
	}

	s.fileAccessMutex.Lock()
	defer s.fileAccessMutex.Unlock()
	data, err := s.readJsonDataFromFile(s.loginDataFilePath)
	if err != nil {
		return err
	}
	found := false
	for index, entry := range data {
		if entry.UserId == userId {
			if data[index].State != OnVacation {
				return errors.New("can not set account to active, when it is not on vacation mode")
			}
			data[index].State = Active
			found = true
			break
		}
	}
	if found {
		err = s.writeJsonDataToFile(s.loginDataFilePath, data)
		if err != nil {
			return err
		}
	} else {
		return errors.New("user not found")
	}

	s.readFileAndUpdateCache(s.loginDataFilePath)
	return nil
}

/*
	Setting your own account to vacation mode.
*/
func (s *LoginSystem) EnableVacationMode(token string) (err error) {
	valid, userId, _, _, err := s.SessionIsValid(token)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("can not set vacation mode for invalid session")
	}

	s.fileAccessMutex.Lock()
	defer s.fileAccessMutex.Unlock()
	data, err := s.readJsonDataFromFile(s.loginDataFilePath)
	if err != nil {
		return err
	}
	found := false
	for index, entry := range data {
		if entry.UserId == userId {
			if data[index].State == WaitingToBeUnlocked {
				return errors.New("can not set a account to vacation mode, when it has not been unlocked")
			}
			data[index].State = OnVacation
			found = true
			break
		}
	}
	if found {
		err = s.writeJsonDataToFile(s.loginDataFilePath, data)
		if err != nil {
			return err
		}
	} else {
		return errors.New("user not found")
	}

	s.readFileAndUpdateCache(s.loginDataFilePath)
	return nil
}

/*
	Get all users which are locked.
*/
func (s *LoginSystem) GetAllLockedUsers() []User {
	s.cachedUserDataMutex.RLock()
	defer s.cachedUserDataMutex.RUnlock()

	var lockedUsers []User
	for _, storedUser := range s.cachedUserData {
		if storedUser.State == WaitingToBeUnlocked {
			user := User{Mail: storedUser.Mail, UserId: storedUser.UserId,
				FirstName: storedUser.FirstName, LastName: storedUser.LastName,
				Role: storedUser.Role, State: storedUser.State}
			lockedUsers = append(lockedUsers, user.Copy())
		}
	}

	return lockedUsers
}

/*
	Get all users which are active.
*/
func (s *LoginSystem) GetAllActiveUsers() []User {
	s.cachedUserDataMutex.RLock()
	defer s.cachedUserDataMutex.RUnlock()

	var lockedUsers []User
	for _, storedUser := range s.cachedUserData {
		if storedUser.State == Active {
			user := User{Mail: storedUser.Mail, UserId: storedUser.UserId,
				FirstName: storedUser.FirstName, LastName: storedUser.LastName,
				Role: storedUser.Role, State: storedUser.State}
			lockedUsers = append(lockedUsers, user.Copy())
		}
	}

	return lockedUsers
}

/*
	Unlock a account which is waiting to be unlocked. The current session needs the permission to do this.
*/
func (s *LoginSystem) UnlockAccount(currentToken string, userIdToUnlock int) (unlocked bool, err error) {
	valid, userId, _, _, err := s.SessionIsValid(currentToken)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("current session is not valid")
	}

	isAdmin := s.checkIfUserIsInAdminRole(userId)
	if !isAdmin {
		return false, errors.New("current session has no permission to unlock accounts")
	}

	s.fileAccessMutex.Lock()
	defer s.fileAccessMutex.Unlock()

	data, err := s.readJsonDataFromFile(s.loginDataFilePath)
	if err != nil {
		return false, err
	}
	found := false
	for index, entry := range data {
		if entry.UserId == userIdToUnlock {
			if data[index].State != WaitingToBeUnlocked {
				return false, errors.New("can not unlock a account, which is not in the waiting to be unlocked state")
			}
			data[index].State = Active
			found = true
			break
		}
	}
	if found {
		err = s.writeJsonDataToFile(s.loginDataFilePath, data)
		if err != nil {
			return false, err
		}
	} else {
		return false, errors.New("user to unlock not found")
	}

	s.readFileAndUpdateCache(s.loginDataFilePath)
	return true, nil
}

/*
	Initialize the files for the user system.
*/
func (s *LoginSystem) initializeFiles(folderPath string) (err error) {
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

/*
	Create a session for a user. Returns the token or an error.
*/
func (s *LoginSystem) createSessionForUser(user storedUserData) (authToken string, err error) {
	s.currentSessionsMutex.Lock()
	defer s.currentSessionsMutex.Unlock()
	token, err := helpers.GenerateUUID()
	if err != nil {
		return "", err
	}
	s.currentSessions[token] = inMemorySession{userMail: user.Mail, userId: user.UserId,
		userRole:         user.Role,
		sessionToken:     authToken,
		sessionTimestamp: time.Now()}
	return token, nil
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
			if strings.ToLower(v.Mail) == strings.ToLower(userName) {
				return true
			}
		}
	}
	return false
}

/*
	Register a new user.
*/
func (s *LoginSystem) registerNewUser(userName string, password string, firstName string, lastName string, role UserRole, state UserState) (err error) {
	s.fileAccessMutex.Lock()
	defer s.fileAccessMutex.Unlock()
	existingData, err := s.readJsonDataFromFile(s.loginDataFilePath)
	if err != nil {
		return err
	}
	generatedLoginData, err := s.generateLoginData(userName, password, len(existingData)+1, firstName, lastName, role, state)
	if err != nil {
		return err
	}
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
	Check if the provided password for a user is correct.
*/
func (s *LoginSystem) checkUserCredentials(user storedUserData, providedPassword string) (success bool) {
	// Build the combination of the provided password with the stored hash and compare it with the stored password hash:
	comb := string(user.StoredSalt) + string(providedPassword)
	passwordHash := sha512.New()
	_, err := io.WriteString(passwordHash, comb)
	if err != nil {
		return false
	}
	correctComb := bytes.Equal(user.StoredPass, passwordHash.Sum(nil))

	return correctComb
}

/*
	Generates the login data.
*/
func (s *LoginSystem) generateLoginData(userName string, password string, newId int,
	firstName string, lastName string, role UserRole, state UserState) (userData storedUserData, err error) {

	// Hashing and salt generation is from: https://www.socketloop.com/tutorials/golang-securing-password-with-salt
	// In a real life scenario we would use something like https://godoc.org/golang.org/x/crypto/bcrypt
	// But we should only use standard packages, so we use this solution.

	// Generate a salt for combination with password.
	salt, err := generateSalt([]byte(password))
	if err != nil {
		return storedUserData{}, err
	}
	// Generate a hash for the password, but combine it with the salt.
	passwordHash := sha512.New()
	combination := string(salt) + string(password)
	_, err = io.WriteString(passwordHash, combination)
	if err != nil {
		return storedUserData{}, err
	}
	// Store the hashed password together with the salt.
	return storedUserData{Mail: userName,
		StoredPass: passwordHash.Sum(nil),
		StoredSalt: salt,
		UserId:     newId,
		FirstName:  firstName,
		LastName:   lastName,
		Role:       role,
		State:      state}, nil
}

/*
	Generate a salt out of a password.
*/
func generateSalt(secret []byte) ([]byte, error) {
	buf := make([]byte, saltSize, saltSize+sha512.Size)
	_, err := io.ReadFull(rand.Reader, buf)

	if err != nil {
		return []byte{}, err
	}

	hash := sha512.New()
	hash.Write(buf)
	hash.Write(secret)
	return hash.Sum(buf), nil
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

	s.cachedUserDataMutex.Lock()
	defer s.cachedUserDataMutex.Unlock()
	parsedData, err := s.readJsonDataFromFile(filePath)
	if parsedData != nil {
		s.cachedUserData = parsedData
	} else {
		s.cachedUserData = *new([]storedUserData)
	}
	return
}

/*
	Check if the userData for the given id has the admin role.
*/
func (s *LoginSystem) checkIfUserIsInAdminRole(userId int) (isAdmin bool) {
	s.cachedUserDataMutex.RLock()
	defer s.cachedUserDataMutex.RUnlock()
	for _, entry := range s.cachedUserData {
		if entry.UserId == userId {
			if entry.Role == Admin {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (s *LoginSystem) changeUserPassword(user User, newPassword string) (bool, error) {
	s.fileAccessMutex.Lock()
	defer s.fileAccessMutex.Unlock()

	data, err := s.readJsonDataFromFile(s.loginDataFilePath)
	if err != nil {
		return false, err
	}
	found := false
	for index, entry := range data {
		if entry.UserId == user.UserId {
			generatedLoginData, err := s.generateLoginData(user.Mail, newPassword, user.UserId, user.FirstName,
				user.LastName, user.Role, user.State)
			if err != nil {
				return false, err
			}
			data[index] = generatedLoginData
			found = true
			break
		}
	}
	if found {
		err = s.writeJsonDataToFile(s.loginDataFilePath, data)
		if err != nil {
			return false, err
		}
	} else {
		return false, errors.New("user to unlock not found")
	}

	s.readFileAndUpdateCache(s.loginDataFilePath)
	return true, nil
}

const saltSize = 32
