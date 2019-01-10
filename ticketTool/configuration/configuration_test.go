// 5894619, 6720876, 9793350
package configuration

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
	Example for validation of a configuration with a invalid unacknowledged mail path.
*/
func ExampleConfiguration_ValidateConfiguration_UnacknowledMailPathIsInvalid_ReturnsFalse() {
	logger := logging.ConsoleLogger{}
	config := Configuration{}
	config.BindFlags()

	config.UnAcknowledgedMailPath = ""
	valid := config.ValidateConfiguration(logger)
	fmt.Print(valid)
	//Output:
	// <Error>[Configuration]: Unacknowledged mails file path must be set.
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
func ExampleConfiguration_ValidateConfiguration_CertificateDoesNotExist_ReturnsFalse() {
	logger := logging.ConsoleLogger{}
	config := Configuration{}
	config.BindFlags()
	config.UnAcknowledgedMailPath = "test.json"

	// The path should not exist.
	config.CertificatePath = "/Temp/test/jspemdusoem.key"
	valid := config.ValidateConfiguration(logger)
	fmt.Println(valid)
	//Output:
	// <Error>[Configuration]: Certificate path does not exist
	// false
}

/*
	Validation should fail, if the certificate file does not exist.
*/
func ExampleConfiguration_ValidateConfiguration_ApiKeyFileDoesNotExist_ReturnsFalse() {
	logger := logging.ConsoleLogger{}
	config := Configuration{}
	config.BindFlags()
	config.UnAcknowledgedMailPath = "test.json"

	// Create a temporary file for the certificate. Otherwise the api key file would not be checked.
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	// The path should not exist.
	config.CertificatePath = tmpfile.Name()
	config.ApiKeysFilePath = "/Temp/test/jspemdusoem.key"
	valid := config.ValidateConfiguration(logger)
	fmt.Println(valid)
	//Output:
	// <Error>[Configuration]: Api keys file path does not exist
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
