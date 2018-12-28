package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
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

func (c *ApiClient) buildPostRequest(url string, data []byte) (*http.Request, error) {
	// TODO: API KEY AFTER Merge: https://stackoverflow.com/questions/12756782/go-http-post-and-use-cookies
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *ApiClient) SendMails(mails []mail.Mail) error{
	jsonData, err := json.Marshal(mails)
	if err != nil {
		return err
	}
	url := "https://" +  c.baseUrl +  ":" + strconv.Itoa(c.port) + sendPath
	req , err := c.buildPostRequest(url, jsonData)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Returned status code: " + strconv.Itoa(resp.StatusCode))
	}
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

   return apiClient, nil
}