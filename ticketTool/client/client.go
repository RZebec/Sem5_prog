package client

import (
	"crypto/tls"
	"crypto/x509"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"io/ioutil"
	"net/http"
)

type Client interface {
	SendMails(mails []mail.Mail) error
	ReceiveMails() ([]mail.Mail, error)
	AcknowledgeMails(mailIds []string) error
}

const sendPath = "/api/mail/incoming"
const receivePath = "/api/mail/outgoing"

type ApiClient struct {
	baseUrl string
	port int
	client *http.Client
}


func (c *ApiClient) SendMails(mails []mail.Mail) error{
	return nil
}

func (c *ApiClient) ReceiveMails() ([]mail.Mail, error){
	return *new([]mail.Mail), nil
}

func (c *ApiClient) AcknowledgeMails(mailIds []string) error {
	return nil
}

func CreateClient(config configuration.Configuration) (ApiClient, error){
	apiClient := ApiClient{}
	apiClient.baseUrl = config.BaseUrl
	apiClient.port = config.Port

	caCert, err := ioutil.ReadFile(config.CertificatePath)
   if err != nil {
			  return apiClient, err
	  }
   caCertPool := x509.NewCertPool()
   caCertPool.AppendCertsFromPEM(caCert)

   apiClient.client = &http.Client{
	   Transport: &http.Transport{ TLSClientConfig: &tls.Config{ RootCAs: caCertPool }, },}

}