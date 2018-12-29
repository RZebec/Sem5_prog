package configuration

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"flag"
	"github.com/pkg/errors"
	"strconv"
)

type Configuration struct {
	Port            int
	BaseUrl         string
	CertificatePath string
	ApiKeysFilePath string
}

func (c *Configuration) RegisterFlags() {
	flag.StringVar(&c.BaseUrl, "baseUrl", "localhost", "the base url")
	flag.IntVar(&c.Port, "port", 9000, "the port to use")
	flag.StringVar(&c.CertificatePath, "certificatePath", "cert.pem", "path to the certificate")
	flag.StringVar(&c.ApiKeysFilePath, "apiKeysFilePAth", "api.keys", "path to the api key file")

}

/*
	Bind the flags to the configuration.
*/
func (c *Configuration) BindFlags() {
	flag.Parse()
}

/*
	Validate the configuration and use the logger to print out the validation error.
*/
func (c *Configuration) ValidateConfiguration(log logging.Logger) bool {
	if c.Port < 0 || c.Port > 65535 {
		log.LogError("Configuration", errors.New("Invalid port. Provided value: "+strconv.Itoa(c.Port)))
		return false
	}

	exists, err := helpers.FilePathExists(c.CertificatePath)
	if err != nil {
		log.LogError("Configuration", errors.Wrap(err, "Could not validate certificate path"))
		return false
	}
	if !exists {
		log.LogError("Configuration", errors.New("Certificate path does not exist"))
		return false
	}

	exists, err = helpers.FilePathExists(c.ApiKeysFilePath)
	if err != nil {
		log.LogError("Configuration", errors.Wrap(err, "Could not validate api key file path"))
		return false
	}
	if !exists {
		log.LogError("Configuration", errors.New("Api keys file path does not exist"))
		return false
	}

	return true
}
