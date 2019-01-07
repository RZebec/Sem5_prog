package ticketData

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Test the conversion from TicketState Open to string.
*/
func TestTicketState_String_Open(t *testing.T) {
	assert.Equal(t, "Open", Open.String())
}

/*
	Test the conversion from TicketState Processing to string.
*/
func TestTicketState_String_Processing(t *testing.T) {
	assert.Equal(t, "Processing", Processing.String())
}

/*
	Test the conversion from TicketState Closed to string.
*/
func TestTicketState_String_Closed(t *testing.T) {
	assert.Equal(t, "Closed", Closed.String())
}
