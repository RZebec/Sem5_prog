package config

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"fmt"
	"testing"
)

func TestConfiguration_GetServiceUrl(t *testing.T) {

}

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

func ExampleConfiguration_ValidateConfiguration_PortIsToBig_ReturnsFalse() {
	logger := logging.ConsoleLogger{}
	config := Configuration{}
	config.BindFlags()

	config.Port = 99999999
	valid := config.ValidateConfiguration(logger)
	fmt.Print(valid)
	//Output:
	// <Error>[Configuration]: Invalid port. Provided value: 99999999
	// false
}

func TestConfiguration_ValidateConfiguration_CertificateDoesNotExist_ReturnsFalse(t *testing.T) {

}

func TestConfiguration_ValidateConfiguration_CertificateKeyDoesNotExist_ReturnsFalse(t *testing.T) {

}