package mailData

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Test for creation of append message mail subject.
*/
func TestBuildAppendMessageNotificationMailSubject(t *testing.T) {
	value := BuildAppendMessageNotificationMailSubject(22)
	assert.Equal(t, "A new message has been appended to your ticket with TicketId<22>", value,
		"The generated string should be correct")
}

/*
	Test for creation of append message mail content.
*/
func TestBuildAppendMessageNotificationMailContent(t *testing.T) {
	generatedString := BuildAppendMessageNotificationMailContent("testReceiver@test.de", "testSender@test.de",
		"This is the test content of the message")
	expectedString := "Hello testReceiver@test.de,\n testSender@test.de added a message to your ticket.\n " +
		"Content of the message: \n This is the test content of the message"
	assert.Equal(t, expectedString, generatedString,
		"The generated string should be correct")
}

/*
	Test for the unlock notification subject.
*/
func TestBuildUnlockUserNotificationMailSubject(t *testing.T) {
	generatedString := BuildUnlockUserNotificationMailSubject()
	assert.Equal(t, "Your account has been unlocked.", generatedString,
		"The unlock notification subject should be correct.")
}

/*
	Test for unlock notification content.
*/
func TestBuildUnlockUserNotificationMailContent(t *testing.T) {
	generatedString := BuildUnlockUserNotificationMailContent("testReceiver@test.de")
	expectedString := "Hello testReceiver@test.de,\n Your account has been activated by the administrator"
	assert.Equal(t, expectedString, generatedString,
		"The unlock notification content should be correct.")
}

/*
	Test for ticket merge notification mail subject.
*/
func TestBuildTicketMergeNotificationMailSubject(t *testing.T) {
	generatedString := BuildTicketMergeNotificationMailSubject(4, 5)
	expectedString := "Your ticket with TicketId<4> has been merged. New Ticket: TicketId<5>"
	assert.Equal(t, expectedString, generatedString,
		"The ticket merge notification subject should be correct.")
}

/*
	Test for ticket merge notification mail content.
*/
func TestBuildTicketMergeNotificationMailContent(t *testing.T) {
	generatedString := BuildTicketMergeNotificationMailContent("testReceiver@test.de", 5, 9)
	expectedString := "Hello testReceiver@test.de,\n Ticket 5 has been merged with 9."
	assert.Equal(t, expectedString, generatedString,
		"The ticket merge notification content should be correct.")
}

/*
	Test for the editor changed notification mail subject.
*/
func TestBuildTicketEditorChangedNotificationMailSubject(t *testing.T) {
	generatedString := BuildTicketEditorChangedNotificationMailSubject(22)
	expectedString := "Your ticket with TicketId<22> has been changed. A new editor has been set or removed"
	assert.Equal(t, expectedString, generatedString,
		"The ticket editor changed notification subject should be correct.")
}

/*
	Test for the editor changed notification mail content.
*/
func TestBuildTicketEditorChangedNotificationMailContent(t *testing.T) {
	generatedString := BuildTicketEditorChangedNotificationMailContent("testReceiver@test.de", 23, "NewEditor")
	expectedString := "Hello testReceiver@test.de,\n Ticket 23 has a new Editor: NewEditor."
	assert.Equal(t, expectedString, generatedString,
		"The ticket editor changed notification content should be correct.")
}

/*
	Test for the editor removed mail notification.
*/
func TestBuildTicketEditorRemovedNotificationMailContent(t *testing.T) {
	generatedString := BuildTicketEditorRemovedNotificationMailContent("testReceiver@test.de", 23)
	expectedString := "Hello testReceiver@test.de,\n The editor has been removed from Ticket 23."
	assert.Equal(t, expectedString, generatedString,
		"The ticket editor removed notification content should be correct.")
}

/*
	Test for the ticket state change notification subject generation.
*/
func TestBuildTicketStateChangedNotificationMailSubject(t *testing.T) {
	generatedString := BuildTicketStateChangedNotificationMailSubject(23)
	expectedString := "The state of your ticket with TicketId<23> has been changed."
	assert.Equal(t, expectedString, generatedString,
		"The ticket state changed notification subject should be correct.")
}

/*
	Test for the ticket state change notification content generation.
*/
func TestBuildTicketStateChangedNotificationMailContent(t *testing.T) {
	generatedString := BuildTicketStateChangedNotificationMailContent("testReceiver@test.de", ticketData.Open)
	expectedString := "Hello testReceiver@test.de,\n The state of your ticket has been changed to Open."
	assert.Equal(t, expectedString, generatedString,
		"The ticket state changed notification content should be correct.")
}
