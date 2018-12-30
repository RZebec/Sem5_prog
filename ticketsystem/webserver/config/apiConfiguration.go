package config

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"encoding/json"
	"errors"
	"sync"
)

/*
	The Api Configuration provides the configuration for the incoming and outgoing mail api key.
*/
type IApiConfiguration interface {
	// Get the incoming mail api key.
	GetIncomingMailApiKey() string
	// Get the outgoing mail api key.
	GetOutgoingMailApiKey() string
	// Change the incoming mail api key.
	ChangeIncomingMailApiKey(newKey string) error
	// Change the outgoing mail api key.
	ChangeOutgoingMailApiKey(newKey string) error
	// Validate the configuration.
	Validate() (bool, string)
}

/*
	The configuration for the api.
*/
type ApiConfiguration struct {
	incomingMailApiKey string
	outgoingMailApiKey string
	filePath           string
	accessMutex        sync.RWMutex
	Test               string
}

/*
	A type to persist data.
*/
type persistedData struct {
	IncomingMailApiKey string
	OutgoingMailApiKey string
}

/*
	Get the incoming mail api key.
*/
func (c *ApiConfiguration) GetIncomingMailApiKey() string {
	c.accessMutex.RLock()
	defer c.accessMutex.RUnlock()
	return c.incomingMailApiKey
}

/*
	Get the outgoing mail api key.
*/
func (c *ApiConfiguration) GetOutgoingMailApiKey() string {
	c.accessMutex.RLock()
	defer c.accessMutex.RUnlock()
	return c.outgoingMailApiKey
}

/*
	Change the incoming mail api key.
*/
func (c *ApiConfiguration) ChangeIncomingMailApiKey(newKey string) error {
	oldValue := c.incomingMailApiKey
	c.incomingMailApiKey = newKey
	valid, msg := c.Validate()
	if !valid {
		c.incomingMailApiKey = oldValue
		return errors.New(msg)
	} else {
		err := c.persist()
		if err != nil {
			return err
		}
		return nil
	}
}

/*
	Change the outgoing mail api key.
*/
func (c *ApiConfiguration) ChangeOutgoingMailApiKey(newKey string) error {
	oldValue := c.outgoingMailApiKey
	c.outgoingMailApiKey = newKey
	valid, msg := c.Validate()
	if !valid {
		c.outgoingMailApiKey = oldValue
		return errors.New(msg)
	} else {
		err := c.persist()
		if err != nil {
			return err
		}
		return nil
	}
}

/*
	Create and initialize the configuration.
*/
func CreateAndInitialize(config Configuration) (*ApiConfiguration, error) {
	existed, err := helpers.CreateFileWithPathIfNotExists(config.ApiKeyFilePath)
	if err != nil {
		return nil, err
	}
	if !existed {
		// Create default api keys:
		config := ApiConfiguration{
			incomingMailApiKey: "MW0j6HXw0QrksRz0lcKUisoJqkAjAhcgs1MFjFWIfTUWoccWmYhBippVGzD5I8dVyx6GXdpkTbOONeAuGw1HreDWbswBMGpnx9Lrk7rglPfaWqLzguAMJdnX7PFhOhbj",
			outgoingMailApiKey: "VKoUaBZtBJA6ZFy7bmyJif2YaGvkPnE6c9tkwYhbGs4cFPcoY2Brv8cdoAEW4Eer0x3OBMtXhNlSBtnKWXKNR6J5MYU8gYo6vrPkSdwhmKxj0su7dYCCSpSeyDShPuY3",
			filePath:           config.ApiKeyFilePath}
		err = config.persist()
		if err != nil {
			return nil, err
		}
		return &config, nil
	}
	apiConfig, err := readKeysFromFile(config.ApiKeyFilePath)
	if err != nil {
		return nil, err
	}
	valid, msg := apiConfig.Validate()
	if !valid {
		return nil, errors.New(msg)
	}
	apiConfig.filePath = config.ApiKeyFilePath
	return &apiConfig, nil
}

/*
	Validate the configuration.
*/
func (c *ApiConfiguration) Validate() (bool, string) {
	if c.incomingMailApiKey == "" || len(c.incomingMailApiKey) < 128 {
		return false, "Incoming mail api key needs to have at least 128 characters"
	}
	if c.outgoingMailApiKey == "" || len(c.outgoingMailApiKey) < 128 {
		return false, "Outgoing mail api key needs to have at least 128 characters"
	}
	return true, ""
}

/*
	Persist the configuration.
*/
func (c *ApiConfiguration) persist() error {
	c.accessMutex.Lock()
	defer c.accessMutex.Unlock()
	jsonData, err := json.MarshalIndent(c.transformToPersistenceData(), "", "    ")
	if err != nil {
		return err
	}
	return helpers.WriteDataToFile(c.filePath, jsonData)

}

/*
	Transform the configuration to a persistence format.
*/
func (c *ApiConfiguration) transformToPersistenceData() interface{} {
	return persistedData{IncomingMailApiKey: c.incomingMailApiKey, OutgoingMailApiKey: c.outgoingMailApiKey}
}

/*
	Read the api keys from a given file.
*/
func readKeysFromFile(filePath string) (ApiConfiguration, error) {
	fileValue, err := helpers.ReadAllDataFromFile(filePath)
	if err != nil {
		return ApiConfiguration{}, err
	}
	parsedData := new(persistedData)
	err = json.Unmarshal(fileValue, &parsedData)
	if err != nil {
		return ApiConfiguration{}, err
	}
	return ApiConfiguration{incomingMailApiKey: parsedData.IncomingMailApiKey, outgoingMailApiKey: parsedData.OutgoingMailApiKey}, nil
}
