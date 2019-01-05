package ticketData

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	A copy of a message should not be able to change the orignal message.
*/
func TestMessageEntry_Copy(t *testing.T) {
	originalMessage := MessageEntry{Id: 1, CreatorMail: "test@test.de", Content: "This is a test", OnlyInternal: false}
	copied := originalMessage.Copy()

	// Change the copy:
	copied.Id = 200
	copied.Content = "changed text"
	copied.CreatorMail = "changed@changed.de"
	copied.OnlyInternal = true

	// Ensure, that the original message has not been changed
	assert.Equal(t, 1, originalMessage.Id, "Original id should not be changed")
	assert.Equal(t, "This is a test", originalMessage.Content, "Original content should not be changed")
	assert.Equal(t, "test@test.de", originalMessage.CreatorMail, "Original creator mail should not be changed")
	assert.Equal(t, false, originalMessage.OnlyInternal, "Original only internal flag should not be changed")
}
