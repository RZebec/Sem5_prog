package mockedForTests

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"github.com/stretchr/testify/mock"
)

/*
	A mocked ticketData context.
*/
type MockedTicketContext struct {
	mock.Mock
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) GetTicketsForEditorId(userId int) []ticketData.TicketInfo {
	args := m.Called(userId)
	return args.Get(0).([]ticketData.TicketInfo)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) GetTicketsForCreatorMail(mail string) []ticketData.TicketInfo {
	args := m.Called(mail)
	return args.Get(0).([]ticketData.TicketInfo)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) CreateNewTicketForInternalUser(title string, editor userData.User, initialMessage ticketData.MessageEntry) (*ticketData.Ticket, error) {
	args := m.Called(title, editor, initialMessage)
	return args.Get(0).(*ticketData.Ticket), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) CreateNewTicket(title string, creator ticketData.Creator, initialMessage ticketData.MessageEntry) (*ticketData.Ticket, error) {
	args := m.Called(title, creator, initialMessage)
	return args.Get(0).(*ticketData.Ticket), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) GetTicketById(id int) (bool, *ticketData.Ticket) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*ticketData.Ticket)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) GetAllTicketInfo() []ticketData.TicketInfo {
	args := m.Called()
	return args.Get(0).([]ticketData.TicketInfo)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) AppendMessageToTicket(ticketId int, message ticketData.MessageEntry) (*ticketData.Ticket, error) {
	args := m.Called(ticketId, message)
	return args.Get(0).(*ticketData.Ticket), args.Error(1)
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
func (m *MockedTicketContext) SetEditor(editor userData.User, ticketId int) (*ticketData.Ticket, error) {
	args := m.Called(editor, ticketId)
	return args.Get(0).(*ticketData.Ticket), args.Error(1)
}

/*
	Mocked function.
*/
func (m *MockedTicketContext) SetTicketState(ticketId int, newState ticketData.TicketState) (*ticketData.Ticket, error) {
	args := m.Called(ticketId, newState)
	return args.Get(0).(*ticketData.Ticket), args.Error(1)
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
func (m *MockedTicketContext) GetAllOpenTickets() []ticketData.TicketInfo {
	args := m.Called()
	return args.Get(0).([]ticketData.TicketInfo)
}
