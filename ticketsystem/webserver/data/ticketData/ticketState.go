// 5894619, 6720876, 9793350
package ticketData

type TicketState int

/*
	The ticket states.
*/
const (
	Open TicketState = 1 + iota
	Processing
	Closed
)

/*
	The states as strings.
*/
var states = [...]string{
	"Open",
	"Processing",
	"Closed",
}

/*
	For the string representation.
*/
func (state TicketState) String() string {
	return states[state-1]
}
