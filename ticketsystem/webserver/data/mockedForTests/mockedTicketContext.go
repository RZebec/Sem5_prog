package mockedForTests

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"github.com/stretchr/testify/mock"
)

/*
	A mocked ticket context.
*/
type MockedTicketContext struct {
	mock.Mock
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) CreateNewTicketForInternalUser(title string, editor user.User, initialMessage ticket.MessageEntry) (*ticket.Ticket, error) {
	args := m.Called(title, editor, initialMessage)
	return args.Get(0).(*ticket.Ticket), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) CreateNewTicket(title string, creator ticket.Creator, initialMessage ticket.MessageEntry) (*ticket.Ticket, error) {
	args := m.Called(title, creator, initialMessage)
	return args.Get(0).(*ticket.Ticket), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) GetTicketById(id int) (bool, *ticket.Ticket) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*ticket.Ticket)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) GetAllTicketInfo() []ticket.TicketInfo {
	args := m.Called()
	return args.Get(0).([]ticket.TicketInfo)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) AppendMessageToTicket(ticketId int, message ticket.MessageEntry) (*ticket.Ticket, error) {
	args := m.Called(ticketId, message)
	return args.Get(0).(*ticket.Ticket), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) MergeTickets(firstTicketId int, secondTicketId int) (success bool, err error) {
	args := m.Called(firstTicketId, secondTicketId)
	return args.Bool(0), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) SetEditor(editor user.User, ticketId int) (*ticket.Ticket, error) {
	args := m.Called(editor, ticketId)
	return args.Get(0).(*ticket.Ticket), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) SetTicketState(ticketId int, newState ticket.TicketState) (*ticket.Ticket, error) {
	args := m.Called(ticketId, newState)
	return args.Get(0).(*ticket.Ticket), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) RemoveEditor(ticketId int) error {
	args := m.Called(ticketId)
	return args.Error(0)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) GetAllOpenTickets() []ticket.TicketInfo {
	args := m.Called()
	return args.Get(0).([]ticket.TicketInfo)
}
