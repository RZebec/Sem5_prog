package ticket

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

func TestTicketManager_GetAllTicketInfo_NoTickets_EmptyArray(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetAllTicketInfo()
	assert.Equal(t, 0, len(tickets))
}

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

func readTicketFromFile(filePath string) (*Ticket, error) {
	ticket, err := initializeFromFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not load ticker")
	}
	return ticket, nil
}

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

	assertNewTickedWithStoredTicket(t, newTicket, storedTicket)
}

func TestTicketManager_CreateNewTicketForInternalUser_TicketCreated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	user := user.User{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 45, CreatorMail: user.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicketForInternalUser("newTestTitle", user, initialMessage)

	exists, createdTicket := testee.GetTicketById(newTicket.info.Id)
	assert.True(t, exists, "the ticket should be created")

	// Validate that the file is created:
	exists, _ = helpers.FilePathExists(createdTicket.filePath)
	assert.True(t, exists, "the ticket file should be created")

	storedTicket, err := readTicketFromFile(createdTicket.filePath)

	assertNewTickedWithStoredTicket(t, newTicket, storedTicket)
}

func assertNewTickedWithStoredTicket(t *testing.T, expected *Ticket, actual *Ticket) {
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

func TestTicketManager_Initialize_NoTicketsExist_Initialized(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	assert.Equal(t, 0, len(testee.cachedTickets))
	assert.Equal(t, 0, len(testee.cachedTicketIds))
}

func TestTicketManager_AppendMessageToTicket(t *testing.T) {
	assert.True(t, false, "not implemented")
	// TODO: Add test for tickets
	// TODO: Add test for wrong initialization
	// TODO: Add multiple tickets in a concurrent way
	// TODO: Add examples
}

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
