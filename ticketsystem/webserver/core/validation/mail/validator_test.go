package mail

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Example for the mail validator.
 */
func ExampleValidator_Validate() {
	val := New()

	mail := "max@mustermann.de"
	fmt.Println(val.Validate(mail))
	// Output:
	// true
}

/*
	Test a invalid mail.
 */
func TestValidator_Validate_InvalidMail_NotOk(t *testing.T) {
	val := New()

	mail := "@mustermann.de"
	assert.False(t, val.Validate(mail), "The mail should be invalid")

	mail = "max@.de"
	assert.False(t, val.Validate(mail), "The mail should be invalid")

	mail = "max@mustermannde"
	assert.False(t, val.Validate(mail), "The mail should be invalid")

	mail = "max@mustermann."
	assert.False(t, val.Validate(mail), "The mail should be invalid")

	mail = "max@mustermann"
	assert.False(t, val.Validate(mail), "The mail should be invalid")
}
