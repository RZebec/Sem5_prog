package ticketData

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
	A copy of a ticket should not be able to change the original ticket.
*/
func TestTicketInfo_Copy(t *testing.T) {
	refTimestamp := time.Now()
	origUser := userData.User{Mail: "test@test.de", UserId: 23, FirstName: "Hans", LastName: "Müller"}
	creator := Creator{Mail: origUser.Mail, FirstName: origUser.FirstName, LastName: origUser.LastName}
	origTicketInfo := TicketInfo{Id: 5, Title: "OrigTitle", HasEditor: true,
		CreationTime: refTimestamp, LastModificationTime: refTimestamp,
		Editor: origUser, Creator: creator, State: Open}

	copiedTicket := origTicketInfo.Copy()

	// Change the copy:
	copiedTicket.Creator.LastName = "changed"
	copiedTicket.Editor.FirstName = "changed"
	copiedTicket.HasEditor = false
	copiedTicket.Title = "changed"
	copiedTicket.CreationTime = time.Now()
	copiedTicket.LastModificationTime = time.Now()
	copiedTicket.State = Closed

	// Assert that the original ticket info has not been changed:
	assert.Equal(t, "Müller", origTicketInfo.Creator.LastName, "Original creator name should not be changed")
	assert.Equal(t, "Hans", origTicketInfo.Editor.FirstName, "Original editor name should not be changed")
	assert.Equal(t, true, origTicketInfo.HasEditor, "Original has editor flag should not be changed")
	assert.Equal(t, "OrigTitle", origTicketInfo.Title, "Original title should not be changed")
	assert.Equal(t, refTimestamp, origTicketInfo.CreationTime, "Original creation time should not be changed")
	assert.Equal(t, refTimestamp, origTicketInfo.LastModificationTime, "Original last modification time should not be changed")
	assert.Equal(t, Open, origTicketInfo.State, "Original state should not be changed")
}
