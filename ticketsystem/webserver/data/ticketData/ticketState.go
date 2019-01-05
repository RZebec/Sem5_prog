package ticketData

type TicketState int

/*
	The ticketData states.
 */
const (
	Open TicketState = 1 + iota
	Processing
	Closed
)

/*
	The states as strings.
 */
var states = [...]string {
	"Open",
	"Processing",
	"Closed",
}

/*
	For the string representation.
 */
func (state TicketState) String() string {
	return states[state -1]
}