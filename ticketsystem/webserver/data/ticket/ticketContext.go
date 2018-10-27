package ticket

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
	"regexp"
	"sync"
	"time"
)

/*
	Interface for the ticket context.
*/
type TicketContext interface {
	CreateNewTicketForInternalUser(title string, editor user.User, initialMessage MessageEntry) (*Ticket, error)
	CreateNewTicket(title string, creator Creator, initialMessage MessageEntry) (*Ticket, error)
	GetTicketById(id int) (Ticket, error)
	GetAllTicketInfo() []TicketInfo
	AppendMessageToTicket(ticketId int, message MessageEntry) (*Ticket, error)
}

// TODO: Set LastModification Time and Create method for info modification
/*
	Append a message to a ticket.
*/
func (t *TicketManager) AppendMessageToTicket(ticketId int, message MessageEntry) (*Ticket, error) {
	exists, ticket := t.GetTicketById(ticketId)
	if exists {
		ticket.messages = append(ticket.messages, message)
		err := ticket.persist()
		if err != nil {
			return nil, errors.Wrap(err, "could not append message to ticket")
		}
		return ticket, nil
	} else {
		return nil, errors.New("ticket does not exist")
	}
}

/*
	Get all existing ticket information.
*/
func (t *TicketManager) GetAllTicketInfo() []TicketInfo {
	t.cachedTicketsMutex.RLock()
	defer t.cachedTicketsMutex.RUnlock()

	var ticketInfos []TicketInfo
	for _, ticket := range t.cachedTickets {
		ticketInfos = append(ticketInfos, ticket.info)
	}
	return ticketInfos
}

/*
	Get a ticket by its id. Returns false if the ticket does not exist.
*/
func (t *TicketManager) GetTicketById(id int) (bool, *Ticket) {
	value, ok := t.cachedTickets[id]
	if ok {
		return true, &value
	}
	return false, nil
}

/*
	The ticket manager handles the access to the tickets.
*/
type TicketManager struct {
	cachedTickets      map[int]Ticket
	cachedTicketIds    []int
	cachedTicketsMutex sync.RWMutex
	ticketFolderPath   string
}

/*
	Initialize the TicketManager with the given folder path. The folder is used to load and store the ticket data.
*/
func (t *TicketManager) Initialize(folderPath string) error {
	if folderPath == "" {
		return errors.New("path to login data storage can not be a empty string.")
	}
	t.cachedTickets = make(map[int]Ticket)
	t.cachedTicketIds = []int{}
	t.ticketFolderPath = folderPath
	folderExists, err := helpers.FilePathExists(folderPath)
	if err != nil {
		return errors.Wrap(err, "could not check if folder for ticket exists.")
	}
	if !folderExists {
		err = helpers.CreateFolderPath(folderPath)
		if err != nil {
			return errors.Wrap(err, "could not create folder for tickets.")
		}
	} else {
		return t.readExistingTickets()
	}
	return nil
}

/*
	Create a new ticket for a internal user. The
*/
func (t *TicketManager) CreateNewTicketForInternalUser(title string, editor user.User, initialMessage MessageEntry) (*Ticket, error) {
	t.cachedTicketsMutex.Lock()
	defer t.cachedTicketsMutex.Unlock()

	newId := maxIntInArray(t.cachedTicketIds) + 1
	t.cachedTicketIds = append(t.cachedTicketIds, newId)

	newTicket, err := createNewEmptyTicket(t.ticketFolderPath, newId)
	if err != nil {
		return nil, errors.Wrap(err, "could not create ticket")
	}
	newTicket.info.Id = newId
	newTicket.info.Creator = ConvertToCreator(editor)
	newTicket.info.HasEditor = true
	newTicket.info.Editor = editor
	newTicket.info.Title = title
	initialMessage.Id = 0
	initialMessage.CreatorMail = editor.Mail
	initialMessage.CreationTime = time.Now()
	newTicket.messages = append(newTicket.messages, initialMessage)

	newTicket.persist()
	t.cachedTickets[newId] = *newTicket
	t.cachedTicketIds = append(t.cachedTicketIds, newId)
	return newTicket, nil

}

/*
	Create a new ticket.
*/
func (t *TicketManager) CreateNewTicket(title string, creator Creator, initialMessage MessageEntry) (*Ticket, error) {

	t.cachedTicketsMutex.Lock()
	defer t.cachedTicketsMutex.Unlock()

	newId := maxIntInArray(t.cachedTicketIds) + 1

	newTicket, err := createNewEmptyTicket(t.ticketFolderPath, newId)
	if err != nil {
		return nil, errors.Wrap(err, "could not create ticket")
	}
	newTicket.info.Id = newId
	newTicket.info.Creator = creator
	newTicket.info.HasEditor = false
	newTicket.info.Title = title
	initialMessage.Id = 0
	initialMessage.CreatorMail = creator.Mail
	initialMessage.CreationTime = time.Now()
	newTicket.messages = append(newTicket.messages, initialMessage)

	newTicket.persist()
	t.cachedTickets[newId] = *newTicket
	t.cachedTicketIds = append(t.cachedTicketIds, newId)
	return newTicket, nil
}

/*
	Read the existing tickets from the file system.
*/
func (t *TicketManager) readExistingTickets() error {
	t.cachedTicketsMutex.Lock()
	defer t.cachedTicketsMutex.Unlock()

	files, err := ioutil.ReadDir(t.ticketFolderPath)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		// Ignore folders
		if !fileInfo.IsDir() {
			match, _ := regexp.MatchString(".json", fileInfo.Name())
			if match {
				ticket, err := initializeFromFile(path.Join(t.ticketFolderPath, fileInfo.Name()))
				if err != nil {
					return errors.Wrap(err, "could not load ticker")
				}
				t.cachedTickets[ticket.info.Id] = *ticket
				t.cachedTicketIds = append(t.cachedTicketIds, ticket.info.Id)
			}
		}
	}
	return nil
}

/*
	Get the max value in an array.
*/
func maxIntInArray(values []int) int {
	if len(values) == 0 {
		return 0
	}
	max := 0
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
