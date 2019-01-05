package config

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/validation/mail"
	"flag"
	"github.com/pkg/errors"
	"strconv"
)

/*
	A struct to hold configuration values.
*/
type Configuration struct {
	LoginDataFolderPath   string
	TicketDataFolderPath  string
	MailDataFolderPath    string
	Port                  int
	BaseUrl               string
	CertificatePath       string
	CertificateKeyPath    string
	ApiKeyFilePath        string
	SendingMailAddress    string
	AccessTokenCookieName string
}

/*
	Register the configuration properties as flags.
*/
func (c *Configuration) RegisterFlags() {
	flag.StringVar(&c.BaseUrl, "baseUrl", "localhost", "the base url")
	flag.IntVar(&c.Port, "port", 9000, "the port to use")
	flag.StringVar(&c.LoginDataFolderPath, "loginDataFolderPath", "data/login", "path to the folder to store the login data")
	flag.StringVar(&c.TicketDataFolderPath, "ticketDataFolderPath", "data/tickets", "path to the folder to store the ticketData data")
	flag.StringVar(&c.MailDataFolderPath, "mailDataFolderPath", "data/mails", "path to the folder to store the mail data")
	flag.StringVar(&c.CertificateKeyPath, "certificateKeyPath", "key.pem", "path to the certificate key file")
	flag.StringVar(&c.CertificatePath, "certificatePath", "cert.pem", "path to the certificate")
	flag.StringVar(&c.ApiKeyFilePath, "apiKeysFilePath", "data/api.keys", "path to the apiKey file")
	flag.StringVar(&c.SendingMailAddress, "sendingMailAddress", "notification@ticketsystem.de", "Outgoing mail address for notifications")
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

	exists, err = helpers.FilePathExists(c.CertificateKeyPath)
	if err != nil {
		log.LogError("Configuration", errors.Wrap(err, "Could not validate certificate key path"))
		return false
	}
	if !exists {
		log.LogError("Configuration", errors.New("Certificate key path does not exist"))
		return false
	}

	mailValidator := mail.NewValidator()
	if !mailValidator.Validate(c.SendingMailAddress) {
		log.LogError("Configuration", errors.New("Outgoing mail address is not valid"))
		return false
	}

	return true
}

/*
	Get the service url. Consisting of the base url and the port.
*/
func (c Configuration) GetServiceUrl() string {
	return c.BaseUrl + ":" + strconv.Itoa(c.Port)
}
