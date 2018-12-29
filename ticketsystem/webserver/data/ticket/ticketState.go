package ticket

type TicketState int

const (
	Open TicketState = 1 + iota
	Processing
	Closed
)
