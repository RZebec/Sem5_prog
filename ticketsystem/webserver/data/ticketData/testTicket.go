package ticketData

/*
	Creates a ticketData for test purposes.
 */
func CreateTestTicket(ticketInfo TicketInfo, messages []MessageEntry) *Ticket {
	ticket := new(Ticket)
	ticket.info = ticketInfo
	ticket.messages = messages
	ticket.filePath = ""

	return ticket
}
