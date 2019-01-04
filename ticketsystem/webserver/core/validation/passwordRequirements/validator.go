package passwordRequirements

import (
	"unicode"
)

/*
	A validator for password requirements.
*/
type Validator struct {
}

/*
	Create a new validator.
*/
func NewValidator() *Validator {
	validator := Validator{}
	return &validator
}

/*
	Validate a password against the requirements. Returns true if valid, false if not.
*/
func (v Validator) Validate(value string) bool {
	number, upperCase, lowerCase, specialCharacter := false, false, false, false

	if len([]rune(value)) < 8 {
		return false
	}
	for _, s := range value {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upperCase = true
		case unicode.IsLower(s):
			lowerCase = true
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			specialCharacter = true
		case unicode.IsLetter(s) || s == ' ':
		}
	}
	return number && upperCase && lowerCase && specialCharacter
}
