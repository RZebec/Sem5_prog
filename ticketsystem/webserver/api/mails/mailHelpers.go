package mails

import (
	"regexp"
	"strconv"
	"strings"
)

/*
	The mailIdExtractor
 */
type mailIdExtractor struct {
	regex *regexp.Regexp
}

/*
	Create a new mail id extractor.
 */
func newMailIdExtractor() *mailIdExtractor{
	reg, _ := regexp.Compile(`<\d+>`)
	ext := mailIdExtractor{}
	ext.regex = reg
	return &ext
}

/*
	Get the ticket id from a string.
 */
func (m *mailIdExtractor) getTicketId(text string) (bool, int){
	hasId := m.regex.MatchString(text)
	if hasId {
		value := m.regex.FindStringSubmatch(text)
		if len(value) > 0 {
			stringValue := value[0]
			stringValue = strings.Replace(stringValue, "<", "",1)
			stringValue = strings.Replace(stringValue, ">", "",1)
			inVal, err :=strconv.Atoi(stringValue)

			if err != nil {
				return false, -1
			}
			return true, inVal
		} else {
			return false, -1
		}

	} else {
		return false, -1
	}
}