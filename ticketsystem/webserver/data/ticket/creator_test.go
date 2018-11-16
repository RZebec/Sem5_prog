package ticket

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Test the conversion of a user to creator.
*/
func TestConvertToCreator(t *testing.T) {
	user := user.User{UserId: 1, Mail: "test@test.de", FirstName: "Alex", LastName: "Müller"}
	creator := ConvertToCreator(user)

	assert.Equal(t, user.Mail, creator.Mail)
	assert.Equal(t, user.FirstName, creator.FirstName)
	assert.Equal(t, user.LastName, creator.LastName)
}

/*
	Changing a copy should change the original creator.
*/
func TestCreator_Copy(t *testing.T) {
	origCreator := Creator{Mail: "orig@test.de", FirstName: "Max", LastName: "Müller"}
	copy := origCreator.Copy()

	// Change the copy
	copy.LastName = "changed"
	copy.FirstName = "changed"
	copy.Mail = "1234@test.de"

	// Ensure that the original creator has not been changed
	assert.Equal(t, "orig@test.de", origCreator.Mail)
	assert.Equal(t, "Max", origCreator.FirstName)
	assert.Equal(t, "Müller", origCreator.LastName)
}
