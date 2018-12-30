package mails

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Correct ticket id should be returned.
 */
func TestMailIdExtractor_GetTicketId_IdReturned(t *testing.T){
	testee := newMailIdExtractor()

	containsId, id := testee.getTicketId("Ticket<25>")
	assert.True(t, containsId, "Should contain a ticket id")
	assert.Equal(t, 25, id,"Should return correct ticket id")

	containsId, id = testee.getTicketId("Ticket<1>")
	assert.True(t, containsId, "Should contain a ticket id")
	assert.Equal(t, 1, id,"Should return correct ticket id")

	containsId, id = testee.getTicketId("Ticket<192>")
	assert.True(t, containsId, "Should contain a ticket id")
	assert.Equal(t, 192, id,"Should return correct ticket id")
}

/*
	Strings without ticket ids should return false.
 */
func TestMailIdExtractor_GetTicketId_NoIdReturned(t *testing.T){
	testee := newMailIdExtractor()

	containsId, id := testee.getTicketId("Ticket<>")
	assert.False(t, containsId, "Should not contain a ticket id")
	assert.Equal(t, -1, id,"Should return -1")

	containsId, id = testee.getTicketId("Ticket<1")
	assert.False(t, containsId, "Should not contain a ticket id")
	assert.Equal(t, -1, id,"Should return -1")

	containsId, id = testee.getTicketId("92>")
	assert.False(t, containsId, "Should not contain a ticket id")
	assert.Equal(t, -1, id,"Should return -1")
}