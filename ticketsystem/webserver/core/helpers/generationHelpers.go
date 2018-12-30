package helpers

import (
	"crypto/rand"
	"fmt"
)

/*
	Generate a UUID.
	Source: https://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language
*/
func GenerateUUID() (string, error) {
	b := make([]byte, 16)
	_, er := rand.Read(b)
	if er != nil {
		return "", er
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
