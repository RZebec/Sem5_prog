// 5894619, 6720876, 9793350
package ticketData

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
	A copied ticket should not be able to change the original ticket.
*/
func TestTicket_Copy(t *testing.T) {
	refTimestamp := time.Now()
	origUser := userData.User{Mail: "test@test.de", UserId: 23, FirstName: "Hans", LastName: "Müller"}
	creator := Creator{Mail: origUser.Mail, FirstName: origUser.FirstName, LastName: origUser.LastName}
	origTicketInfo := TicketInfo{Id: 5, Title: "OrigTitle", HasEditor: true,
		CreationTime: refTimestamp, LastModificationTime: refTimestamp,
		Editor: origUser, Creator: creator}
	originalMessage := MessageEntry{Id: 1, CreatorMail: "test@test.de", Content: "This is a test", OnlyInternal: false}
	origMessages := []MessageEntry{originalMessage}

	origTicket := Ticket{info: origTicketInfo, messages: origMessages}

	copied := origTicket.Copy()

	// Change the copy
	copied.info.Title = "changed"
	copied.messages[0].Content = "changed"

	// Ensure that the original ticket has not been changed:
	assert.Equal(t, "OrigTitle", origTicket.Info().Title, "The original info should not be changed")
	assert.Equal(t, "This is a test", origTicket.Messages()[0].Content,
		"The original info should not be changed")

}
