package mail

import "regexp"

/*
	A validator for mail addresses.
*/
type Validator struct {
	regex *regexp.Regexp
}

/*
	Create a new validator.
*/
func New() *Validator {
	// Regex from: https://stackoverflow.com/questions/23968992/how-to-match-a-regex-with-backreference-in-go
	reg, _ := regexp.Compile(`^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	validator := Validator{}
	validator.regex = reg
	return &validator
}

/*
	Validate a mail address. Returns true if valid, false if not.
*/
func (v Validator) Validate(value string) bool {
	return v.regex.MatchString(value)
}
