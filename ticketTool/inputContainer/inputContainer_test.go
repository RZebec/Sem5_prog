package inputContainer

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Example for validation of a configuration with a invalid (negative) port value.
*/
func TestExampleConfiguration_ValidateConfiguration_PortIsNegative_ReturnsFalse(t *testing.T) {
	logger := logging.ConsoleLogger{}
	config := Configuration{}
	config.BindFlags()

	config.Port = -5
	valid := config.ValidateConfiguration(logger)
	fmt.Print(valid)
	//Output:
	// <Error>[Configuration]: Invalid port. Provided value: -5
	// false
}

/*
	Example for validation of a configuration with a invalid (too big) port value.
*/
func TestExampleConfiguration_ValidateConfiguration_PortIsToBig_ReturnsFalse(t *testing.T) {
	logger := logging.ConsoleLogger{}
	config := Configuration{}
	config.BindFlags()

	config.Port = 99999999
	valid := config.ValidateConfiguration(logger)
	fmt.Println(valid)
	//Output:
	// <Error>[Configuration]: Invalid port. Provided value: 99999999
	// false
}

/*
	Validation should fail, if the certificate file does not exist.
*/
func TestConfiguration_ValidateConfiguration_CertificateDoesNotExist_ReturnsFalse(t *testing.T) {
	logger := logging.ConsoleLogger{}
	config := Configuration{}
	config.BindFlags()

	// The path should not exist.
	config.CertificatePath = "/Temp/test/jspemdusoem.key"
	valid := config.ValidateConfiguration(logger)
	fmt.Println(valid)
	//Output:
	// <Error>[Configuration]: Certificate path does not exist
	// false
}

/*
	Default values should be set.
*/
func TestConfiguration_RegisterAndBindFlags(t *testing.T) {
	config := Configuration{}

	config.RegisterFlags()
	config.BindFlags()

	// Default values should be set:
	assert.Equal(t, "localhost", config.BaseUrl, "default value for host should be set")
	assert.Equal(t, 9000, config.Port, "default value for port should be set")
	assert.Equal(t, "cert.pem", config.CertificatePath, "default value for certificate should be set")
}
