package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	UserName should be returned.
*/
func TestUser_GetUserNameString(t *testing.T) {
	user := User{UserId: 1, Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"}
	assert.Equal(t, "Max Mustermann", user.GetUserNameString())
}

/*
	A change on a copy should not be able to change the original.
*/
func TestUser_Copy(t *testing.T) {
	user := User{UserId: 1, Mail: "test@test.de", FirstName: "Max", LastName: "Mustermann"}
	copiedUser := user.Copy()
	copiedUser.LastName = "Müller"

	// The original user should not be changed:
	assert.Equal(t, "Mustermann", user.LastName)
}
