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
