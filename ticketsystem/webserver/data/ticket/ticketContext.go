package ticket

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

/*
	Interface for the ticket context.
*/
type TicketContext interface {
	CreateNewTicketForInternalUser(title string, editor user.User, initialMessage MessageEntry) (*Ticket, error)
	CreateNewTicket(title string, creator Creator, initialMessage MessageEntry) (*Ticket, error)
	GetTicketById(id int) (bool, *Ticket)
	GetTicketsForEditorId(userId int) []TicketInfo
	GetTicketsForCreatorMail(mail string) []TicketInfo
	GetAllTicketInfo() []TicketInfo
	AppendMessageToTicket(ticketId int, message MessageEntry) (*Ticket, error)
	MergeTickets(firstTicketId int, secondTicketId int) (success bool, err error)
	SetEditor(editor user.User, ticketId int) (*Ticket, error)
	SetTicketState(ticketId int, newState TicketState) (*Ticket, error)
	RemoveEditor(ticketId int) error
	GetAllOpenTickets() []TicketInfo
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
	Get the ticket for a given creator mail.
 */
func (t *TicketManager) GetTicketsForCreatorMail(mail string) []TicketInfo {
	t.cachedTicketsMutex.RLock()
	defer t.cachedTicketsMutex.RUnlock()

	var ticketInfos []TicketInfo
	for _, ticket := range t.cachedTickets {
		if strings.ToLower(ticket.info.Creator.Mail) == strings.ToLower(mail) {
			ticketInfos = append(ticketInfos, ticket.info.Copy())
		}
	}
	return ticketInfos
}

/*
	Get the ticket for a given editor.
 */
func (t *TicketManager) GetTicketsForEditorId(userId int) []TicketInfo {
	t.cachedTicketsMutex.RLock()
	defer t.cachedTicketsMutex.RUnlock()

	var ticketInfos []TicketInfo
	for _, ticket := range t.cachedTickets {
		if ticket.info.HasEditor && ticket.info.Editor.UserId == userId {
			ticketInfos = append(ticketInfos, ticket.info.Copy())
		}
	}
	return ticketInfos
}

/*
	Get all tickets with the state open.
*/
func (t *TicketManager) GetAllOpenTickets() []TicketInfo {
	t.cachedTicketsMutex.RLock()
	defer t.cachedTicketsMutex.RUnlock()

	var ticketInfos []TicketInfo
	for _, ticket := range t.cachedTickets {
		if ticket.info.State == Open {
			ticketInfos = append(ticketInfos, ticket.info.Copy())
		}
	}
	return ticketInfos
}

/*
	Set the state of a ticket.
*/
func (t *TicketManager) SetTicketState(ticketId int, newState TicketState) (*Ticket, error) {
	exists, ticket := t.GetTicketById(ticketId)
	if exists {
		ticket.info.State = newState
		ticket.info.LastModificationTime = time.Now()
		err := ticket.persist()
		if err != nil {
			return nil, errors.Wrap(err, "could not change state of the ticket")
		}

		t.cachedTicketsMutex.Lock()
		defer t.cachedTicketsMutex.Unlock()
		t.cachedTickets[ticket.info.Id] = *ticket

		return ticket.Copy(), nil
	} else {
		return nil, errors.New("ticket does not exist")
	}
}

/*
	Append a message to a ticket.
*/
func (t *TicketManager) AppendMessageToTicket(ticketId int, message MessageEntry) (*Ticket, error) {
	exists, ticket := t.GetTicketById(ticketId)
	if exists {
		ticket.messages = append(ticket.messages, message)
		// Fix the id of the message
		for i := range ticket.messages {
			ticket.messages[i].Id = i
		}
		ticket.info.LastModificationTime = time.Now()
		err := ticket.persist()
		if err != nil {
			return nil, errors.Wrap(err, "could not append message to ticket")
		}

		t.cachedTicketsMutex.Lock()
		defer t.cachedTicketsMutex.Unlock()
		t.cachedTickets[ticket.info.Id] = *ticket

		return ticket.Copy(), nil
	} else {
		return nil, errors.New("ticket does not exist")
	}
}

/*
	Remove the editor of a ticket.
*/
func (t *TicketManager) RemoveEditor(ticketId int) error {
	exists, ticket := t.GetTicketById(ticketId)
	if exists {
		ticket.info.Editor = user.GetInvalidDefaultUser()
		ticket.info.HasEditor = false

		ticket.info.LastModificationTime = time.Now()
		err := ticket.persist()
		if err != nil {
			return errors.Wrap(err, "could not remove editor of ticket")
		}

		t.cachedTicketsMutex.Lock()
		defer t.cachedTicketsMutex.Unlock()
		t.cachedTickets[ticket.info.Id] = *ticket

		return nil
	} else {
		return errors.New("ticket does not exist")
	}
}

/*
	Set the editor of a ticket.
*/
func (t *TicketManager) SetEditor(editor user.User, ticketId int) (*Ticket, error) {
	exists, ticket := t.GetTicketById(ticketId)
	if exists {
		ticket.info.Editor = editor
		ticket.info.HasEditor = true

		ticket.info.LastModificationTime = time.Now()
		err := ticket.persist()
		if err != nil {
			return nil, errors.Wrap(err, "could not set editor of ticket")
		}

		t.cachedTicketsMutex.Lock()
		defer t.cachedTicketsMutex.Unlock()
		t.cachedTickets[ticket.info.Id] = *ticket

		return ticket.Copy(), nil
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
		ticketInfos = append(ticketInfos, ticket.info.Copy())
	}
	return ticketInfos
}

/*
	Get a ticket by its id. Returns false if the ticket does not exist.
*/
func (t *TicketManager) GetTicketById(id int) (bool, *Ticket) {
	value, ok := t.cachedTickets[id]
	if ok {
		return true, value.Copy()
	}
	return false, nil
}

func (t *TicketManager) MergeTickets(firstTicketId int, secondTicketId int) (success bool, err error) {
	if firstTicketId == secondTicketId {
		return false, errors.New("can not merge a ticket with itself")
	}

	firstTicket, firstTicketFound := t.cachedTickets[firstTicketId]
	secondTicket, secondTicketFound := t.cachedTickets[secondTicketId]
	if !firstTicketFound || !secondTicketFound {
		return false, errors.New("ticket not found")
	}

	if !firstTicket.info.HasEditor || !secondTicket.info.HasEditor {
		return false, errors.New("can not merge ticket if there is no editor")
	}

	if firstTicket.info.Editor.UserId != secondTicket.info.Editor.UserId {
		return false, errors.New("only tickets of the same editor can be merged")
	}

	t.cachedTicketsMutex.Lock()
	defer t.cachedTicketsMutex.Unlock()

	// Merge the message entries:
	olderTicket := firstTicket
	newerTicket := secondTicket
	if secondTicket.info.Id < firstTicket.info.Id {
		olderTicket = secondTicket
		newerTicket = firstTicket
	}

	olderTicket.messages = append(olderTicket.messages, newerTicket.messages...)
	// Fix the id of the message
	for i := range olderTicket.messages {
		olderTicket.messages[i].Id = i
	}

	sort.Slice(olderTicket.messages, func(i, j int) bool {
		return olderTicket.messages[i].CreationTime.Before(olderTicket.messages[j].CreationTime)
	})

	olderTicket.info.LastModificationTime = time.Now()
	olderTicket.persist()
	t.cachedTickets[olderTicket.info.Id] = olderTicket

	// Remove the newer ticket:
	t.cachedTicketIds = remove(t.cachedTicketIds, newerTicket.info.Id)
	filePathToDelete := newerTicket.filePath
	delete(t.cachedTickets, newerTicket.info.Id)
	err = os.Remove(filePathToDelete)

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

/*
	Initialize the TicketManager with the given folder path. The folder is used to load and store the ticket data.
*/
func (t *TicketManager) Initialize(folderPath string) error {
	if folderPath == "" {
		return errors.New("path to ticket data storage can not be a empty string.")
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
	return newTicket.Copy(), nil

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
	return newTicket.Copy(), nil
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
	Remove an int from an int array.
*/
func remove(s []int, i int) []int {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
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
