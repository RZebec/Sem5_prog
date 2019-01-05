package mails

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Correct ticketData id should be returned.
*/
func TestMailIdExtractor_GetTicketId_IdReturned(t *testing.T) {
	testee := newMailIdExtractor()

	containsId, id := testee.getTicketId("Ticket<25>")
	assert.True(t, containsId, "Should contain a ticketData id")
	assert.Equal(t, 25, id, "Should return correct ticketData id")

	containsId, id = testee.getTicketId("Ticket<1>")
	assert.True(t, containsId, "Should contain a ticketData id")
	assert.Equal(t, 1, id, "Should return correct ticketData id")

	containsId, id = testee.getTicketId("Ticket<192>")
	assert.True(t, containsId, "Should contain a ticketData id")
	assert.Equal(t, 192, id, "Should return correct ticketData id")
}

/*
	Strings without ticketData ids should return false.
*/
func TestMailIdExtractor_GetTicketId_NoIdReturned(t *testing.T) {
	testee := newMailIdExtractor()

	containsId, id := testee.getTicketId("Ticket<>")
	assert.False(t, containsId, "Should not contain a ticketData id")
	assert.Equal(t, -1, id, "Should return -1")

	containsId, id = testee.getTicketId("Ticket<1")
	assert.False(t, containsId, "Should not contain a ticketData id")
	assert.Equal(t, -1, id, "Should return -1")

	containsId, id = testee.getTicketId("92>")
	assert.False(t, containsId, "Should not contain a ticketData id")
	assert.Equal(t, -1, id, "Should return -1")
}
