// 5894619, 6720876, 9793350
package acknowledgementStorage

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"encoding/json"
	"github.com/pkg/errors"
	"sync"
)

/*
	Interface for the acknowledgment storage.
*/
type AckStorage interface {
	DeleteAcknowledges(delete []mailData.Acknowledgment) error
	AppendAcknowledgements(acknowledge []mailData.Acknowledgment) error
	ReadAcknowledgements() ([]mailData.Acknowledgment, error)
}

/*
	The AckManager.
*/
type AckManager struct {
	filePath        string
	fileAccessMutex sync.RWMutex
	acknowledgments []mailData.Acknowledgment
}

/*
	Initialize method for the ack manager.
*/
func InitializeAckManager(filepath string) (*AckManager, error) {
	if filepath == "" {
		return nil, errors.New("filepath for storage of acknowledgements is invalid")
	}

	manager := AckManager{filePath: filepath}
	// Create the file if it does not exist
	existed, err := helpers.CreateFileWithPathIfNotExists(filepath)
	if err != nil {
		return nil, err
	}
	if existed {
		err := manager.readDataFromFile()
		if err != nil {
			return nil, err
		}
	} else {
		err = manager.writeDataToFile()
		if err != nil {
			return nil, err
		}
	}
	return &manager, nil
}

/*
	Delete acknowledgements from the tracked list.
*/
func (m *AckManager) DeleteAcknowledges(delete []mailData.Acknowledgment) error {
	var acksToSave []mailData.Acknowledgment
	for _, currentAck := range m.acknowledgments {
		deleteAck := false
		for _, toDelete := range delete {
			if currentAck.Id == toDelete.Id {
				deleteAck = true
				break
			}
		}
		if !deleteAck {
			acksToSave = append(acksToSave, currentAck)
		}
	}

	m.acknowledgments = acksToSave
	return m.writeDataToFile()
}

/*
	Write new Acknowledgements to the storage. They will be appended to the existing one.
*/
func (m *AckManager) AppendAcknowledgements(acknowledge []mailData.Acknowledgment) error {
	for _, ackToAdd := range acknowledge {
		m.acknowledgments = append(m.acknowledgments, ackToAdd)
	}
	return m.writeDataToFile()
}

/*
	Read Acknowledgment data from file.
*/
func (m *AckManager) ReadAcknowledgements() ([]mailData.Acknowledgment, error) {
	err := m.readDataFromFile()
	if err != nil {
		return []mailData.Acknowledgment{}, err
	}
	return m.acknowledgments, nil
}

/*
	Read ack data from file.
*/
func (m *AckManager) readDataFromFile() error {
	m.fileAccessMutex.Lock()
	defer m.fileAccessMutex.Unlock()
	fileValue, err := helpers.ReadAllDataFromFile(m.filePath)
	if err != nil {
		return err
	}
	parsedData := new([]mailData.Acknowledgment)
	err = json.Unmarshal(fileValue, &parsedData)
	if err != nil {
		return err
	}
	m.acknowledgments = *parsedData
	return nil
}

/*
	Write the ack data to file.
*/
func (m *AckManager) writeDataToFile() error {
	m.fileAccessMutex.Lock()
	defer m.fileAccessMutex.Unlock()
	if len(m.acknowledgments) == 0 {
		m.acknowledgments = make([]mailData.Acknowledgment, 0)
	}
	jsonData, err := json.MarshalIndent(m.acknowledgments, "", "    ")
	if err != nil {
		return err
	}
	return helpers.WriteDataToFile(m.filePath, jsonData)
}
