package passwordRequirements

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Example how to use the validator.
*/
func ExampleValidator_Validate() {
	val := NewValidator()

	pass := "Test!23?"
	fmt.Println(val.Validate(pass))
	// Output:
	// true
}

/*
	A password which is too short, should be invalid. Min length = 8.
*/
func TestValidator_Validate_PasswordTooShort_Invalid(t *testing.T) {
	val := NewValidator()

	// Length = 7
	pass := "ZI(9§sl"
	assert.False(t, val.Validate(pass), "The password should be invalid")
}

/*
	A password without a uppercase is invalid.
*/
func TestValidator_Validate_NoUppercase_Invalid(t *testing.T) {
	val := NewValidator()

	pass := "test45?%67&"
	assert.False(t, val.Validate(pass), "The password should be invalid")
}

/**
A password without a lowercase is invalid.
*/
func TestValidator_Validate_NoLowercase_Invalid(t *testing.T) {
	val := NewValidator()

	pass := "KDI(/)2686DF9§9$"
	assert.False(t, val.Validate(pass), "The password should be invalid")
}

/*
	A password without a number is invalid.
*/
func TestValidator_Validate_NoNumber_Invalid(t *testing.T) {
	val := NewValidator()

	pass := "TPO)(%/SJTOSO"
	assert.False(t, val.Validate(pass), "The password should be invalid")
}

/*
	A password without a special character is invalid.
*/
func TestValidator_Validate_NoSpecialCharacter_Invalid(t *testing.T) {
	val := NewValidator()

	pass := "LosujdOAK5796SK"
	assert.False(t, val.Validate(pass), "The password should be invalid")
}
