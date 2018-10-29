package ticket

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

/*
	Example for the initialization of the ticket manager.
*/
func ExampleTicketManager_Initialize() {
	ticketContext := TicketManager{}
	err := ticketContext.Initialize("pathToTicketFolder")
	fmt.Println(err)
	// Output:
	// <nil>
}

/*
	Example to get a ticket by its id.
*/
func ExampleTicketManager_GetTicketById() {
	ticketContext := TicketManager{}
	initializeWithTempTicketForExample(&ticketContext)

	exists, ticket := ticketContext.GetTicketById(1)
	fmt.Println(exists)
	fmt.Println(ticket.Info().Id)
	// Output:
	// true
	// 1
}

/*
	Example to get all ticket infos.
*/
func ExampleTicketManager_GetAllTicketInfo() {
	ticketContext := TicketManager{}
	initializeWithTempTicketForExample(&ticketContext)

	ticket := ticketContext.GetAllTicketInfo()
	fmt.Println(ticket[0].Id)
	fmt.Println(ticket[0].Title)
	// Output:
	// 1
	// Example_TicketTitle
}

/*
	Example to create a new ticket for an internal user.
*/
func ExampleTicketManager_CreateNewTicketForInternalUser() {
	// Preparation for example:
	folderPath, rootPath, _ := prepareTempDirectory()
	defer os.RemoveAll(rootPath)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	user := user.User{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: user.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicketForInternalUser("Example_TestTitle", user, initialMessage)

	fmt.Println(err)
	fmt.Println(newTicket.Info().Title)
	// Output:
	// <nil>
	// Example_TestTitle
}

/*
	Example to create a new ticket.
*/
func ExampleTicketManager_CreateNewTicket() {
	// Preparation for example:
	folderPath, rootPath, _ := prepareTempDirectory()
	defer os.RemoveAll(rootPath)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	creator := Creator{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicket("Example_TestTitle", creator, initialMessage)

	fmt.Println(err)
	fmt.Println(newTicket.Info().Title)
	// Output:
	// <nil>
	// Example_TestTitle
}

/*
	Example to append a new message to a ticket.
*/
func ExampleTicketManager_AppendMessageToTicket() {
	ticketContext := TicketManager{}
	folderPath := initializeWithTicketForExample(&ticketContext)
	defer os.RemoveAll(folderPath)

	exists, ticket := ticketContext.GetTicketById(1)
	fmt.Println(exists)
	fmt.Println(len(ticket.Messages()))

	newMessage := MessageEntry{CreatorMail: "alex@wagner.de", Content: "This is the message", OnlyInternal: false}

	updatedTicket, _ := ticketContext.AppendMessageToTicket(ticket.Info().Id, newMessage)

	fmt.Println(len(updatedTicket.Messages()))
	fmt.Println(updatedTicket.Messages()[1].Content)

	// Output:
	// true
	// 1
	// 2
	// This is the message
}

/*
	Getting the ticket infos when no tickets exists, should return a empty array.
*/
func TestTicketManager_GetAllTicketInfo_NoTickets_EmptyArrayReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetAllTicketInfo()
	assert.Equal(t, 0, len(tickets))
}

/*
	Getting the ticket infos when ticket exist, should return the infos for the existing tickets.
*/
func TestTicketManager_GetAllTicketInfo_TicketsExist_TicketInfoReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetAllTicketInfo()
	assert.Equal(t, 4, len(tickets))

	// Assert, that the correct data is returned:
	for _, v := range tickets {
		// Compare ticket 4:
		if v.Id == 4 {
			assert.Equal(t, "TestTitle4", v.Title)
			assert.Equal(t, "peter@test.de", v.Editor.Mail)
			assert.Equal(t, 2, v.Editor.UserId)
			assert.Equal(t, "Peter", v.Editor.FirstName)
			assert.Equal(t, "Test", v.Editor.LastName)
			assert.Equal(t, true, v.HasEditor)
			assert.Equal(t, "peter@test.de", v.Creator.Mail)
			assert.Equal(t, "Peter", v.Creator.FirstName)
			assert.Equal(t, "Test", v.Creator.LastName)
			expectedCreationTime, _ := time.Parse(time.RFC3339Nano, "2018-10-27T18:23:04.1141357+02:00")
			assert.Equal(t, expectedCreationTime, v.CreationTime)
			expectedLastModificationTime, _ := time.Parse(time.RFC3339Nano, "2018-10-27T18:23:04.1141357+02:00")
			assert.Equal(t, expectedLastModificationTime, v.LastModificationTime)
		}
	}
}

/*
	Creating tickets for internal users in a concurrent way, should create all tickets.
*/
func TestTicketManager_CreateNewTicketsForInternalUser_ConcurrentAccess_AllCreated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)
	numberOfCreatedTickets := 100

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(numberOfCreatedTickets)
	for i := 0; i < numberOfCreatedTickets; i++ {

		go func(number int) {
			defer waitGroup.Done()
			id := strconv.Itoa(number)
			user := user.User{Mail: id + "@web.de", UserId: number, FirstName: "firstName" + id, LastName: "lastName" + id}
			initialMessage := MessageEntry{Id: 45, CreatorMail: user.Mail, Content: "This is a test" + id, OnlyInternal: false}
			testee.CreateNewTicketForInternalUser("newTestTitle"+id, user, initialMessage)

		}(i)
	}
	waitGroup.Wait()
	tickets := testee.GetAllTicketInfo()

	assert.Equal(t, numberOfCreatedTickets, len(tickets), "all ticket info should be cached")
	// Check if ticket data is correct
	for i := 0; i < numberOfCreatedTickets; i++ {
		id := strconv.Itoa(i)
		expectedCreator := Creator{Mail: id + "@web.de", FirstName: "firstName" + id, LastName: "lastName" + id}
		expectedTitle := "newTestTitle" + id
		found := false

		for _, ticket := range tickets {
			if ticket.Title == expectedTitle {
				found = true
				assert.Equal(t, expectedCreator, ticket.Creator, "the creator should be set")
			}
		}

		assert.True(t, found, "ticket has not been found")
	}
}

/*
	Creating tickets for external users in a concurrent way, should create all tickets.
*/
func TestTicketManager_CreateNewTickets_ConcurrentAccess_AllCreated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)
	numberOfCreatedTickets := 100

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(numberOfCreatedTickets)
	for i := 0; i < numberOfCreatedTickets; i++ {

		go func(number int) {
			defer waitGroup.Done()
			id := strconv.Itoa(number)
			creator := Creator{Mail: id + "@web.de", FirstName: "Alex" + id, LastName: "Wagner" + id}
			initialMessage := MessageEntry{Id: 45, CreatorMail: creator.Mail, Content: "This is a test" + id, OnlyInternal: false}
			testee.CreateNewTicket("newTestTitle"+id, creator, initialMessage)

		}(i)
	}
	waitGroup.Wait()
	tickets := testee.GetAllTicketInfo()

	assert.Equal(t, numberOfCreatedTickets, len(tickets), "all ticket info should be cached")
	// Check if ticket data is correct
	for i := 0; i < numberOfCreatedTickets; i++ {
		id := strconv.Itoa(i)
		expectedCreator := Creator{Mail: id + "@web.de", FirstName: "Alex" + id, LastName: "Wagner" + id}
		expectedTitle := "newTestTitle" + id
		found := false

		for _, ticket := range tickets {
			if ticket.Title == expectedTitle {
				found = true
				assert.Equal(t, expectedCreator, ticket.Creator, "the creator should be set")
			}
		}

		assert.True(t, found, "ticket has not been found")
	}
}

/*
	Creating a ticket should really create it.
*/
func TestTicketManager_CreateNewTicket_TicketCreated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	creator := Creator{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 45, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicket("newTestTitle", creator, initialMessage)

	exists, createdTicket := testee.GetTicketById(newTicket.info.Id)
	assert.True(t, exists, "the ticket should be created")

	// Validate that the file is created:
	exists, _ = helpers.FilePathExists(createdTicket.filePath)
	assert.True(t, exists, "the ticket file should be created")

	storedTicket, err := readTicketFromFile(createdTicket.filePath)

	assertTicket(t, newTicket, storedTicket)
}

/*
	Creating a ticket for a internal user should really create the ticket.
*/
func TestTicketManager_CreateNewTicketForInternalUser_TicketCreated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	user := user.User{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: user.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicketForInternalUser("newTestTitle", user, initialMessage)

	exists, createdTicket := testee.GetTicketById(newTicket.info.Id)
	assert.True(t, exists, "the ticket should be created")

	// Validate that the file is created:
	exists, _ = helpers.FilePathExists(createdTicket.filePath)
	assert.True(t, exists, "the ticket file should be created")

	storedTicket, err := readTicketFromFile(createdTicket.filePath)

	assertTicket(t, newTicket, storedTicket)
}

/*
	Initializing the ticket manager with existing tickets, should load the ticket.
*/
func TestTicketManager_Initialize_TicketsExist_TicketsAreLoaded(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	assert.Equal(t, 4, len(testee.cachedTickets))
	assert.Equal(t, 4, len(testee.cachedTicketIds))
}

/*
	Initializing the ticket manager when no ticket exist, should initialize without problems.
*/
func TestTicketManager_Initialize_NoTicketsExist_Initialized(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	assert.Equal(t, 0, len(testee.cachedTickets))
	assert.Equal(t, 0, len(testee.cachedTicketIds))
}

/*
	Initializing the ticket manager with a invalid path, should return an error.
*/
func TestTicketManager_Initialize_InvalidFolderPath(t *testing.T) {
	testee := TicketManager{}
	err := testee.Initialize("")
	assert.Error(t, err, "folderPath is invalid")
}

/*
	Appending a message to a ticket should append the message to the ticket.
*/
func TestTicketManager_AppendMessageToTicket_MessageAppended(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	creator := Creator{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 45, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicket("newTestTitle", creator, initialMessage)

	// Id should be set by the context:
	message := MessageEntry{Id: 9999, CreatorMail: "max@muster.de", CreationTime: time.Now(),
		Content: "This is a appended message", OnlyInternal: false}
	ticket, err := testee.AppendMessageToTicket(newTicket.Info().Id, message)
	assert.Nil(t, err)

	found, storedTicket := testee.GetTicketById(newTicket.info.Id)
	assert.True(t, found)
	assert.Equal(t, ticket, storedTicket)

	for i, message := range storedTicket.messages {
		assert.Equal(t, i, message.Id, "the id should be in order")
	}
}

/*
	Appending a message to a non existing ticket should return a error.
*/
func TestTicketManager_AppendMessageToTicket_TicketDoesNotExist_ErrorReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	message := MessageEntry{Id: 9999, CreatorMail: "max@muster.de", CreationTime: time.Now(),
		Content: "This is a appended message", OnlyInternal: false}
	_, err = testee.AppendMessageToTicket(999, message)
	assert.Equal(t, "ticket does not exist", err.Error())

}

/*
	Write the test data to a folder.
*/
func writeTestDataToFolder(folderPath string) error {
	sampleData := []byte(firstTestTicket)
	// First ticket: id 1
	sampleDataPath := path.Join(folderPath, "1.json")
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err := ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}

	// Second ticket: id 2
	sampleData = []byte(secondTestTicket)
	sampleDataPath = path.Join(folderPath, "2.json")
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err = ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}

	// Third ticket: id 3
	sampleData = []byte(thirdTestTicket)
	sampleDataPath = path.Join(folderPath, "3.json")
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err = ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}

	// Fourth ticket: id 4
	sampleData = []byte(fourthTestTicket)
	sampleDataPath = path.Join(folderPath, "4.json")
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err = ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}
	return nil
}

/*
	Prepare a temporary directory for tests.
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
	Read a ticket from a given file.
*/
func readTicketFromFile(filePath string) (*Ticket, error) {
	ticket, err := initializeFromFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not load ticker")
	}
	return ticket, nil
}

/*
	Asserting a ticket.
*/
func assertTicket(t *testing.T, expected *Ticket, actual *Ticket) {
	// Id is set, so this can not be compared:
	// FilePath can not be compared
	// Creation- and LastModificationTime can not be compared

	// Compare ticket info
	assert.Equal(t, expected.info.Title, actual.info.Title, "the title should be stored")
	assert.Equal(t, expected.info.Editor, actual.info.Editor, "the editor should be stored")
	assert.Equal(t, expected.info.HasEditor, actual.info.HasEditor, "the hasEditor field should be stored")
	assert.Equal(t, expected.info.Creator, actual.info.Creator, "the creator should be stored")

	// Compare messages
	assert.Equal(t, len(expected.messages), len(actual.Messages()))
	sort.Slice(expected.messages, func(i, j int) bool {
		return expected.messages[i].Id < expected.messages[j].Id
	})
	sort.Slice(actual.messages, func(i, j int) bool {
		return actual.messages[i].Id < actual.messages[j].Id
	})

	for i, expectedMessage := range expected.messages {
		actualMessage := actual.messages[i]
		// Creation time can not be compared
		assert.Equal(t, expectedMessage.Id, actualMessage.Id, "message id should be equal")
		assert.Equal(t, expectedMessage.CreatorMail, actualMessage.CreatorMail, "creator should be equal")
		assert.Equal(t, expectedMessage.OnlyInternal, actualMessage.OnlyInternal, " only internal should be equal")
		assert.Equal(t, expectedMessage.Content, actualMessage.Content, "content should be equal")
	}
}
func initializeWithTempTicketForExample(c *TicketManager) {
	folderPath, rootPath, _ := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	c.Initialize(folderPath)

	creator := Creator{Mail: "test@test.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	c.CreateNewTicket("Example_TicketTitle", creator, initialMessage)
}

func initializeWithTicketForExample(c *TicketManager) string {
	folderPath, rootPath, _ := prepareTempDirectory()
	c.Initialize(folderPath)

	creator := Creator{Mail: "test@test.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	c.CreateNewTicket("Example_TicketTitle", creator, initialMessage)
	return rootPath
}
