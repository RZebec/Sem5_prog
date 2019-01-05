package config

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

/*
	Example for validation of a configuration with a invalid (negative) port value.
*/
func ExampleConfiguration_ValidateConfiguration_PortIsNegative_ReturnsFalse() {
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
func ExampleConfiguration_ValidateConfiguration_PortIsToBig_ReturnsFalse() {
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
	Validation should fail, if the outgoing mail address is not valid.
*/
func TestConfiguration_ValidateConfiguration_OutgoingMailAddressInvalid_ReturnsFalse(t *testing.T) {
	// Create a temporary file for the files. Otherwise the mail address would not be checked.
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up
	logger := logging.ConsoleLogger{}
	config := Configuration{}
	config.BindFlags()

	config.CertificatePath = tmpfile.Name()
	config.CertificateKeyPath = tmpfile.Name()
	config.ApiKeyFilePath = tmpfile.Name()
	config.SendingMailAddress = "123@lk"

	valid := config.ValidateConfiguration(logger)

	fmt.Println(valid)
	//Output:
	// <Error>[Configuration]: Outgoing mail address is not valid
	// false
}

/*
	Validation should fail, if the certificate key does not exist.
*/
func ExampleConfiguration_ValidateConfiguration_CertificateKeyDoesNotExist_ReturnsFalse() {
	// Create a temporary file for the certificate. Otherwise the certificate key file would not be checked.
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	logger := logging.ConsoleLogger{}
	config := Configuration{}
	config.BindFlags()

	// The path should not exist.
	config.CertificatePath = tmpfile.Name()
	config.CertificateKeyPath = "Temp/sdsodk"

	valid := config.ValidateConfiguration(logger)
	fmt.Println(valid)
	//Output:
	// <Error>[Configuration]: Certificate key path does not exist
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
	assert.Equal(t, "key.pem", config.CertificateKeyPath, "default value for certificate key file should be set")
	assert.Equal(t, "cert.pem", config.CertificatePath, "default value for certificate should be set")
	assert.Equal(t, "data/tickets", config.TicketDataFolderPath, "default value for ticket data path should be set")
	assert.Equal(t, "data/login", config.LoginDataFolderPath, "default value for login data path should be set")

	assert.Equal(t, "localhost:9000", config.GetServiceUrl(), "url and port should be concatenated")
}
