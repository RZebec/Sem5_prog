// 5894619, 6720876, 9793350
package acknowledgementStorage

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

/*
	The manage should be able to be initialized, even without a existing file.
*/
func TestInitializeAckManager_NoDataExists(t *testing.T) {
	// Create a temporary directory for the storage.
	dir, err := ioutil.TempDir("", "tests")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	filePath := path.Join(dir, "nonExisting.json")

	storage, createErr := InitializeAckManager(filePath)
	assert.Nil(t, createErr)

	assert.Equal(t, 0, len(storage.acknowledgments), "Storage should be created without error")
}

/*
	The data should be loaded, when there is data.
*/
func TestInitializeAckManager_DataExists_DataLoaded(t *testing.T) {
	// Create a temporary directory for the storage.
	dir, err := ioutil.TempDir("", "tests")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	filePath := path.Join(dir, "nonExisting.json")
	writeTestDataToFile(t, filePath)

	testee, createErr := InitializeAckManager(filePath)
	assert.Nil(t, createErr)
	assert.Equal(t, 3, len(testee.acknowledgments), "Three acks should be loaded")
}

/*
	Deleting a acknowledgement should be possible,
*/
func TestAckManager_DeleteAcknowledges(t *testing.T) {
	// Create a temporary directory for the storage.
	dir, err := ioutil.TempDir("", "tests")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	filePath := path.Join(dir, "nonExisting.json")
	writeTestDataToFile(t, filePath)

	testee, createErr := InitializeAckManager(filePath)
	assert.Nil(t, createErr)

	acks, err := testee.ReadAcknowledgements()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(acks), "Three acks should be returned")

	acksToDelete := []mailData.Acknowledgment{acks[0], acks[1]}
	err = testee.DeleteAcknowledges(acksToDelete)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(testee.acknowledgments), "Only one ack should be there")

	persistedData := readPersistedDataFromFile(t, filePath)
	assert.NotNil(t, persistedData)
	assert.Equal(t, 1, len(persistedData), "Only one ack should be persisted")
}

/*
	Appending Acknowledgements should be possible.
*/
func TestAckManager_AppendAcknowledgements(t *testing.T) {
	// Create a temporary directory for the storage..
	dir, err := ioutil.TempDir("", "tests")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	filePath := path.Join(dir, "nonExisting.json")

	testee, createErr := InitializeAckManager(filePath)
	assert.Nil(t, createErr)

	firstAck := mailData.Acknowledgment{Id: "testId01", Subject: "TestSubject01"}
	secondAck := mailData.Acknowledgment{Id: "testId02", Subject: "TestSubject02"}
	acks := []mailData.Acknowledgment{firstAck, secondAck}

	err = testee.AppendAcknowledgements(acks)
	assert.Equal(t, 2, len(testee.acknowledgments), "Both acknowledgments should be stored in memory")

	thirdAck := mailData.Acknowledgment{Id: "testId03", Subject: "TestSubject03"}
	acks = []mailData.Acknowledgment{thirdAck}
	err = testee.AppendAcknowledgements(acks)
	assert.Equal(t, 3, len(testee.acknowledgments), "New acknowledgement should be appended")

	persistedData := readPersistedDataFromFile(t, filePath)
	assert.NotNil(t, persistedData)
	assert.Equal(t, 3, len(persistedData), "All three acks should be persisted")

	for idx, inMemoryAck := range testee.acknowledgments {
		persistedAck := persistedData[idx]
		assert.Equal(t, inMemoryAck, persistedAck, "The persisted data should be equal to the in memory data")
	}

}

/*
	Reading acknowledgements from file should be possible.
*/
func TestAckManager_ReadAcknowledgements(t *testing.T) {
	// Create a temporary directory for the storage.
	dir, err := ioutil.TempDir("", "tests")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	filePath := path.Join(dir, "nonExisting.json")
	writeTestDataToFile(t, filePath)

	testee, createErr := InitializeAckManager(filePath)
	assert.Nil(t, createErr)

	acks, err := testee.ReadAcknowledgements()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(acks), "Three acks should be returned")
}

/*
	Write test data to file.
*/
func writeTestDataToFile(t *testing.T, filePath string) {
	data := `[
    {
        "Id": "testId01",
        "Subject": "TestSubject01"
    },
    {
        "Id": "testId02",
        "Subject": "TestSubject02"
    },
    {
        "Id": "testId03",
        "Subject": "TestSubject03"
    }
]`
	err := helpers.WriteDataToFile(filePath, []byte(data))
	assert.Nil(t, err)
}

/*
	Read persisted data from file.
*/
func readPersistedDataFromFile(t *testing.T, filePath string) []mailData.Acknowledgment {
	fileValue, err := helpers.ReadAllDataFromFile(filePath)
	assert.Nil(t, err)

	parsedData := new([]mailData.Acknowledgment)
	err = json.Unmarshal(fileValue, &parsedData)
	assert.Nil(t, err)

	return *parsedData
}
