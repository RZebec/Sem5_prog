package ticketData

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Test the conversion of a userData to creator.
*/
func TestConvertToCreator(t *testing.T) {
	testUser := userData.User{UserId: 1, Mail: "test@test.de", FirstName: "Alex", LastName: "Müller"}
	creator := ConvertToCreator(testUser)

	assert.Equal(t, testUser.Mail, creator.Mail)
	assert.Equal(t, testUser.FirstName, creator.FirstName)
	assert.Equal(t, testUser.LastName, creator.LastName)
}

/*
	Changing a copy should change the original creator.
*/
func TestCreator_Copy(t *testing.T) {
	origCreator := Creator{Mail: "orig@test.de", FirstName: "Max", LastName: "Müller"}
	copied := origCreator.Copy()

	// Change the copy
	copied.LastName = "changed"
	copied.FirstName = "changed"
	copied.Mail = "1234@test.de"

	// Ensure that the original creator has not been changed
	assert.Equal(t, "orig@test.de", origCreator.Mail)
	assert.Equal(t, "Max", origCreator.FirstName)
	assert.Equal(t, "Müller", origCreator.LastName)
}
