package ticket

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"encoding/json"
	"github.com/pkg/errors"
	"path"
	"strconv"
	"sync"
	"time"
)

/*
	A ticket, containing all information.
*/
type Ticket struct {
	info        TicketInfo
	messages    []MessageEntry
	filePath    string
	accessMutex sync.RWMutex
}

/*
	The struct represents the stored ticket data.
*/
type storedTicket struct {
	Info     TicketInfo
	Messages []MessageEntry
}
/*
	Copy a ticket.
 */
func (t *Ticket) Copy() (*Ticket){
	copiedMessages := *new([]MessageEntry)
	for _, message := range t.messages  {
		copiedMessages = append(copiedMessages, message.Copy())
	}
	return &Ticket{info: t.info.Copy(), messages: copiedMessages, filePath: t.filePath, accessMutex: t.accessMutex}
}

/*
	Transforms the ticket data to a store-able data type.
*/
func (t *Ticket) transformToPersistenceData() storedTicket {
	return storedTicket{Info: t.info, Messages: t.messages}
}

/*
	Loads the data from a store-able data type into the ticket.
*/
func (t *Ticket) loadDataFromPersistenceData(storedData storedTicket) {
	t.info = storedData.Info
	t.messages = storedData.Messages
}

/*
	Get the ticket info.
*/
func (t *Ticket) Info() TicketInfo {
	t.accessMutex.RLock()
	defer t.accessMutex.RUnlock()
	return t.info
}

/*
	Get the ticket entries.
*/
func (t *Ticket) Messages() []MessageEntry {
	t.accessMutex.RLock()
	defer t.accessMutex.RUnlock()
	return t.messages
}

/*
	Initialize a ticket from a given file path.
*/
func initializeFromFile(filePath string) (*Ticket, error) {
	if filePath == "" {
		return nil, errors.New("filePath is invalid")
	}
	// Check if path exists:
	exists, err := helpers.FilePathExists(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not check if file exists")
	}
	if !exists {
		return nil, errors.New("file does not exist")
	}

	fileValue, err := helpers.ReadAllDataFromFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read data from file")
	}

	ticket := new(Ticket)
	ticket.accessMutex.Lock()
	defer ticket.accessMutex.Unlock()

	parsedData := new(storedTicket)
	json.Unmarshal(fileValue, &parsedData)
	ticket.filePath = filePath
	ticket.loadDataFromPersistenceData(*parsedData)
	return ticket, nil
}

/*
	Create a new empty ticket, with a given id.
*/
func createNewEmptyTicket(folderPath string, id int) (*Ticket, error) {
	if folderPath == "" {
		return nil, errors.New("folderPath is invalid")
	}
	// Check if path exists:
	exists, err := helpers.FilePathExists(folderPath)
	if err != nil {
		return nil, errors.Wrap(err, "Could not check if folder exists")
	}
	// Create the parent folder if needed:
	if !exists {
		err = helpers.CreateFolderPath(folderPath)
		if err != nil {
			return nil, errors.Wrap(err, "Could not create parent folder")
		}
	}

	filePath := path.Join(folderPath, strconv.Itoa(id)+".json")
	alreadyExists, err := helpers.FilePathExists(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not check if path already exists")
	}
	if alreadyExists {
		return nil, errors.New("file for this ticket id already exists")
	}
	err = helpers.CreateFileIfNotExists(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not create file for ticket")
	}

	ticket := new(Ticket)
	ticket.info.Id = id
	ticket.info.CreationTime = time.Now()
	ticket.info.LastModificationTime = ticket.info.CreationTime
	ticket.filePath = filePath

	err = ticket.persist()
	if err != nil {
		return nil, errors.Wrap(err, "could not persist created, empty ticker")
	}
	return ticket, nil
}

/*
	Persist the ticket.
*/
func (t *Ticket) persist() error {
	if t.filePath == "" {
		return errors.New("filePath is invalid")
	}
	t.accessMutex.Lock()
	defer t.accessMutex.Unlock()
	jsonData, err := json.MarshalIndent(t.transformToPersistenceData(), "", "    ")
	if err != nil {
		return errors.Wrap(err, "could not save ticket to file")
	}
	return helpers.WriteDataToFile(t.filePath, jsonData)
}
