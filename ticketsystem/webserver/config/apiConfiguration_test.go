// 5894619, 6720876, 9793350
package config

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

/*
	Testing the validation of the configuration.
*/
func TestApiConfiguration_Validate(t *testing.T) {
	testee := ApiConfiguration{}

	testee.outgoingMailApiKey = validOutgoingApiKey
	testee.incomingMailApiKey = "1234"

	// Test the validation of the incoming api key
	valid, msg := testee.Validate()
	assert.False(t, valid, "Configuration should not be valid")
	assert.Equal(t, "Incoming mail api key needs to have at least 128 characters", msg)

	testee.outgoingMailApiKey = "123545"
	testee.incomingMailApiKey = validIncomingApiKey

	// Test the validation of the outgoing api key
	valid, msg = testee.Validate()
	assert.False(t, valid, "Configuration should not be valid")
	assert.Equal(t, "Outgoing mail api key needs to have at least 128 characters", msg)

	// Test for valid configuration
	testee.outgoingMailApiKey = validOutgoingApiKey
	testee.incomingMailApiKey = validIncomingApiKey
	valid, msg = testee.Validate()
	assert.True(t, valid, "Configuration should be valid")
	assert.Equal(t, "", msg)
}

/*
	Get the incoming mail api key.
*/
func TestApiConfiguration_GetIncomingMailApiKey(t *testing.T) {
	testee := ApiConfiguration{}

	testee.incomingMailApiKey = validIncomingApiKey

	value := testee.GetIncomingMailApiKey()
	assert.Equal(t, validIncomingApiKey, value, "The correct api key should be returned")
}

/*
	Get the outgoing mail api key.
*/
func TestApiConfiguration_GetOutgoingMailApiKey(t *testing.T) {
	testee := ApiConfiguration{}

	testee.outgoingMailApiKey = validOutgoingApiKey

	value := testee.GetOutgoingMailApiKey()
	assert.Equal(t, validOutgoingApiKey, value, "The correct api key should be returned")
}

/*
	Persisting a configuration should save it values to a file.
*/
func TestApiConfiguration_PersistConfig(t *testing.T) {
	// Create a temporary file for the api keys.
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	testee := ApiConfiguration{}
	testee.outgoingMailApiKey = validOutgoingApiKey
	testee.incomingMailApiKey = validIncomingApiKey
	testee.filePath = tmpfile.Name()

	// Persist the configuration:
	err = testee.persist()
	assert.Nil(t, err)

	// Check, that the values have been stored:
	persistedConfig, err := readKeysFromFile(testee.filePath)
	assert.Nil(t, err)
	assert.Equal(t, validIncomingApiKey, persistedConfig.GetIncomingMailApiKey(), "Api key should be persisted")
	assert.Equal(t, validOutgoingApiKey, persistedConfig.GetOutgoingMailApiKey(), "Api key should be persisted")
}

/*
	Changing the outgoing mail api key should change the in memory key and the persisted key.
*/
func TestApiConfiguration_ChangeOutgoingMailApiKey(t *testing.T) {
	// Create a temporary file for the api keys.
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	testee := ApiConfiguration{}
	testee.outgoingMailApiKey = validOutgoingApiKey
	testee.incomingMailApiKey = validIncomingApiKey

	testee.filePath = tmpfile.Name()
	err = testee.persist()
	assert.Nil(t, err)

	// Ensure that the old keys are stored:
	persistedConfig, err := readKeysFromFile(testee.filePath)
	assert.Nil(t, err)
	assert.Equal(t, validIncomingApiKey, persistedConfig.GetIncomingMailApiKey(), "Api key should be persisted")
	assert.Equal(t, validOutgoingApiKey, persistedConfig.GetOutgoingMailApiKey(), "Api key should be persisted")

	// Change the key:
	newOutgoingKey := validOutgoingApiKey + "5454"
	testee.ChangeOutgoingMailApiKey(newOutgoingKey)

	// Ensure that the persisted keys and the in memory keys are changed:
	persistedConfig, err = readKeysFromFile(testee.filePath)
	assert.Nil(t, err)
	assert.Equal(t, newOutgoingKey, persistedConfig.GetOutgoingMailApiKey(), "Api key should be persisted")
	assert.Equal(t, newOutgoingKey, testee.GetOutgoingMailApiKey(), "Api key should be changed")
}

/*
	Changing the incoming mail key should change the in memory and the persisted key.
*/
func TestApiConfiguration_ChangeIncomingMailApiKey(t *testing.T) {
	// Create a temporary file for the api keys.
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	testee := ApiConfiguration{}
	testee.outgoingMailApiKey = validOutgoingApiKey
	testee.incomingMailApiKey = validIncomingApiKey

	testee.filePath = tmpfile.Name()
	err = testee.persist()
	assert.Nil(t, err)

	// Ensure that the unchanged keys are stored:
	persistedConfig, err := readKeysFromFile(testee.filePath)
	assert.Nil(t, err)
	assert.Equal(t, validIncomingApiKey, persistedConfig.GetIncomingMailApiKey(), "Api key should be persisted")
	assert.Equal(t, validOutgoingApiKey, persistedConfig.GetOutgoingMailApiKey(), "Api key should be persisted")

	// Change the key:
	newIncomingApiKey := validIncomingApiKey + "5454"
	testee.ChangeIncomingMailApiKey(newIncomingApiKey)

	// Assert that the in memory and the persisted key are changed:
	persistedConfig, err = readKeysFromFile(testee.filePath)
	assert.Nil(t, err)
	assert.Equal(t, newIncomingApiKey, persistedConfig.GetIncomingMailApiKey(), "Api key should be persisted")
	assert.Equal(t, newIncomingApiKey, testee.GetIncomingMailApiKey(), "Api key should be changed")
}

/*
	Changing the key to a invalid key should not be possible. The old key should still be active.
*/
func TestApiConfiguration_ChangeIncomingMailApiKey_InvalidKey(t *testing.T) {
	// Create a temporary file for the api keys.
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	testee := ApiConfiguration{}
	testee.outgoingMailApiKey = validOutgoingApiKey
	testee.incomingMailApiKey = validIncomingApiKey

	testee.filePath = tmpfile.Name()
	err = testee.persist()
	assert.Nil(t, err)

	// Ensure that the old key is persisted and in memory:
	config, err := readKeysFromFile(testee.filePath)
	assert.Nil(t, err)
	assert.Equal(t, validIncomingApiKey, config.GetIncomingMailApiKey(), "Api key should be persisted")
	assert.Equal(t, validIncomingApiKey, testee.GetIncomingMailApiKey(), "Api key should be persisted")

	// Try to change to something invalid:
	err = testee.ChangeIncomingMailApiKey("")

	// Ensure that the old key is still active:
	assert.Equal(t, "Incoming mail api key needs to have at least 128 characters", err.Error())
	assert.Equal(t, validIncomingApiKey, config.GetIncomingMailApiKey(), "Api key should be persisted")
	assert.Equal(t, validIncomingApiKey, testee.GetIncomingMailApiKey(), "Api key should be persisted")
}

/*
	Changing the key to a invalid key should not be possible. The old key should still be active.
*/
func TestApiConfiguration_ChangeOutgoingMailApiKey_InvalidKey(t *testing.T) {
	// Create a temporary file for the api keys.
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	testee := ApiConfiguration{}
	testee.outgoingMailApiKey = validOutgoingApiKey
	testee.incomingMailApiKey = validIncomingApiKey

	testee.filePath = tmpfile.Name()
	err = testee.persist()
	assert.Nil(t, err)

	// Ensure that the old key is persisted and in memory:
	config, err := readKeysFromFile(testee.filePath)
	assert.Nil(t, err)
	assert.Equal(t, validIncomingApiKey, config.GetIncomingMailApiKey(), "Api key should be persisted")
	assert.Equal(t, validIncomingApiKey, testee.GetIncomingMailApiKey(), "Api key should be persisted")

	// Try to change to something invalid:
	err = testee.ChangeOutgoingMailApiKey("")

	// Ensure that the old key is still active:
	assert.Equal(t, "Outgoing mail api key needs to have at least 128 characters", err.Error())
	assert.Equal(t, validOutgoingApiKey, config.GetOutgoingMailApiKey(), "Api key should be persisted")
	assert.Equal(t, validOutgoingApiKey, testee.GetOutgoingMailApiKey(), "Api key should be persisted")
}

/*
	Initializing from a existing key file should be possible.
*/
func TestCreateAndInitialize_ExistingKeysLoaded(t *testing.T) {
	// Create a temporary file for the api keys.
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	// Write a key file to a temporary directory:
	writeTestDataToFile(t, tmpfile.Name())
	configuration := WebServerConfiguration{}
	configuration.ApiKeyFilePath = tmpfile.Name()

	// Initialize from existing key file:
	testee, err := CreateAndInitialize(configuration)
	assert.Nil(t, err)

	assert.Equal(t, validIncomingApiKey, testee.incomingMailApiKey, "The stored data should be loaded")
	assert.Equal(t, validOutgoingApiKey, testee.outgoingMailApiKey, "The stored data should be loaded")
}

/*
	Initializing from a non existing file should set the default keys and create a config file.
*/
func TestCreateAndInitialize_DefaultValuesSet(t *testing.T) {
	// Create a temporary directory for the api keys.
	dir, err := ioutil.TempDir("", "tests")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Ensure that the key file does not exist:
	keyFilePath := path.Join(dir, "nonExisting.key")
	exists, err := helpers.FilePathExists(keyFilePath)
	assert.False(t, exists, "The file should not exist")
	assert.Nil(t, err)

	configuration := WebServerConfiguration{}
	configuration.ApiKeyFilePath = keyFilePath

	// Initialize the configuration:
	testee, err := CreateAndInitialize(configuration)
	assert.Nil(t, err)

	// Assert that the in memory keys are set:
	assert.True(t, len(testee.incomingMailApiKey) >= 128, "Default values should be set")
	assert.True(t, len(testee.outgoingMailApiKey) >= 128, "Default values should be set")

	// Assert that a configuration has been created:
	persistedConfig, err := readKeysFromFile(testee.filePath)
	assert.Nil(t, err)
	assert.True(t, len(persistedConfig.incomingMailApiKey) >= 128, "Default values should be set")
	assert.True(t, len(persistedConfig.outgoingMailApiKey) >= 128, "Default values should be set")

	// Assert that the persisted and in memory configuration are equal:
	assert.Equal(t, testee.GetOutgoingMailApiKey(), persistedConfig.outgoingMailApiKey, "Outgoing api key should be equal")
	assert.Equal(t, testee.GetIncomingMailApiKey(), persistedConfig.incomingMailApiKey, "Incoming api key should be equal")
}

/*
	Write the test data to a file.
*/
func writeTestDataToFile(t *testing.T, filePath string) {
	os.MkdirAll(filepath.Dir(filePath), 0644)
	data := persistedData{
		IncomingMailApiKey: validIncomingApiKey,
		OutgoingMailApiKey: validOutgoingApiKey}
	jsonData, err := json.Marshal(data)
	assert.Nil(t, err)
	sampleData := []byte(jsonData)
	err = ioutil.WriteFile(filePath, sampleData, 0644)
	assert.Nil(t, err)
}

const validIncomingApiKey = "nNIr6vgamoa06F15jlnB98GGT5YY5qk4fvSsB3V8uJD3mpZAKlCMQVeEBf5SHsoXMKUBWfqKYpoj991fB3amzmmes0JaXjiTQRERXEDJsZoinD3bngqz7YjIXdNc6kll"

const validOutgoingApiKey = "zMLky9tCxQ6otKrmB3hyq2q4qnzSntW4hAVziRBuLZBh8aHJ5R7Sut72NPDGfazWDidJ0RewjYWKwKCCaBVSWCSMdafA7BWVOKFO5gBvfEj4VfIPO7cBCC0MiCbq0ZLT"
