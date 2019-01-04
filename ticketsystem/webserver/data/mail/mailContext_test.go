package mail

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

/*
	A logger for tests.
*/
func getTestLogger() logging.Logger {
	return logging.ConsoleLogger{SetTimeStamp: false}
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
	Getting the unsent mails when there are no mails, should return a empty array.
*/
func TestMailManager_GetUnsentMails_NoMailsAvailable(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := MailManager{}

	testee.Initialize(folderPath, "test@test.de", getTestLogger())

	mails, err := testee.GetUnsentMails()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(mails), "There should be no unsent mails")
}

/*
	Creating and getting a new mail should be possible.
*/
func TestMailManager_CreateNewOutgoingMail_GetNewMail(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := MailManager{}

	testee.Initialize(folderPath, "test@test.de", getTestLogger())

	err = testee.CreateNewOutgoingMail("testReceiver@test.de", "TestSubject", "TestContent")
	assert.Nil(t, err)

	mails, err := testee.GetUnsentMails()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(mails), "There should be one unsent mails")

	createdMail := mails[0]
	assert.Equal(t, "testReceiver@test.de", createdMail.Receiver, "The correct receiver should be set")
	assert.Equal(t, "TestSubject", createdMail.Subject, "The correct subject should be set")
	assert.Equal(t, "TestContent", createdMail.Content, "The correct content should be set")
	assert.Equal(t, "test@test.de", createdMail.Sender, "The correct sender should be set")

	persistedMail, err := testee.readMailFromFile(path.Join(folderPath, createdMail.Id+".json"))
	assert.Nil(t, err)
	assert.NotNil(t, persistedMail, "Persisted mail should exist")
	assert.Equal(t, "testReceiver@test.de", persistedMail.Receiver, "The correct receiver should be set")
	assert.Equal(t, "TestSubject", persistedMail.Subject, "The correct subject should be set")
	assert.Equal(t, "TestContent", persistedMail.Content, "The correct content should be set")
	assert.Equal(t, "test@test.de", persistedMail.Sender, "The correct sender should be set")
}

/*
	Acknowledging a retrieved mail should be possible.
*/
func TestMailManager_AcknowledgeMails_MailDeleted(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := MailManager{}

	testee.Initialize(folderPath, "test@test.de", getTestLogger())

	err = testee.CreateNewOutgoingMail("testReceiver@test.de", "TestSubject", "TestContent")
	assert.Nil(t, err)

	mails, err := testee.GetUnsentMails()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(mails), "There should be one unsent mails")

	createdMail := mails[0]
	assert.Equal(t, "testReceiver@test.de", createdMail.Receiver, "The correct receiver should be set")
	assert.Equal(t, "TestSubject", createdMail.Subject, "The correct subject should be set")
	assert.Equal(t, "TestContent", createdMail.Content, "The correct content should be set")
	assert.Equal(t, "test@test.de", createdMail.Sender, "The correct sender should be set")

	mailPath := path.Join(folderPath, createdMail.Id+".json")
	exists, err := helpers.FilePathExists(mailPath)
	assert.Nil(t, err)
	assert.True(t, exists, "Persisted mail should exist")

	acknowledgment := Acknowledgment{Id: createdMail.Id, Subject: createdMail.Subject}
	acks := []Acknowledgment{acknowledgment}

	err = testee.AcknowledgeMails(acks)
	assert.Nil(t, err)

	exists, err = helpers.FilePathExists(mailPath)
	assert.Nil(t, err)
	assert.False(t, exists, "Persisted mail should be deleted")
}

/*
	Existing data should be loaded at the initialization.
*/
func TestMailManager_Initialize_ExistingDataLoaded(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared mails:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := MailManager{}

	err = testee.Initialize(folderPath, "test@test.de", getTestLogger())
	assert.Nil(t, err)

	assert.Equal(t, 3, len(testee.unSentMails), "There should be 3 unsent mails")

	mails, err := testee.GetUnsentMails()
	assert.Equal(t, 3, len(mails), "There should be 3 received mails")
	assert.Equal(t, 3, len(testee.unAcknowledgedMails), "There should be 3 unAcknowledged mails")
}

/*
	Acknowledged mails should be deleted.
*/
func TestMailManager_AcknowledgeMails_AcksAndMailDeleted(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared mails:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := MailManager{}

	err = testee.Initialize(folderPath, "test@test.de", getTestLogger())
	assert.Nil(t, err)

	assert.Equal(t, 3, len(testee.unSentMails), "There should be 3 unsent mails")

	mails, err := testee.GetUnsentMails()
	assert.Equal(t, 3, len(mails), "There should be 3 received mails")
	assert.Equal(t, 3, len(testee.unAcknowledgedMails), "There should be 3 unAcknowledged mails")

	var acks []Acknowledgment
	for _, ackMail := range mails {
		acks = append(acks, Acknowledgment{Id: ackMail.Id, Subject: ackMail.Subject})
	}

	err = testee.AcknowledgeMails(acks)
	assert.Nil(t, err)

	// Assert that all mail files are deleted
	for _, m := range mails {
		mailPath := path.Join(folderPath, m.Id+".json")
		exists, err := helpers.FilePathExists(mailPath)
		assert.Nil(t, err)
		assert.False(t, exists, "Persisted mail should be deleted")
	}

	//
	assert.Equal(t, 0, len(testee.unAcknowledgedMails), "There should be no more unacknowledged mails")
}

/*
	Write the test data to a folder.
*/
func writeTestDataToFolder(folderPath string) error {
	sampleData := []byte(firstTestMail)

	sampleDataPath := path.Join(folderPath, firstTestmailFileName)
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err := ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}

	sampleData = []byte(secondTestMail)
	sampleDataPath = path.Join(folderPath, secondTestmailFileName)
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err = ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}

	sampleData = []byte(thirdTestMail)
	sampleDataPath = path.Join(folderPath, thirdTestmailFileName)
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err = ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}
	return nil
}
