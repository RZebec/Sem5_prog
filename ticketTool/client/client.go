package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
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

type ApiClient struct {
	baseUrl         string
	port            int
	client          *http.Client
	sendingApiKey   string
	receivingApiKey string
}

type persistedData struct {
	IncomingMailApiKey string
	OutgoingMailApiKey string
}

func (c *ApiClient) buildPostRequest(url string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Cookie", shared.AuthenticationCookieName+"="+c.sendingApiKey)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *ApiClient) buildGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", shared.AuthenticationCookieName+"="+c.receivingApiKey)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *ApiClient) SendMails(mails []mail.Mail) error {
	jsonData, err := json.Marshal(mails)
	if err != nil {
		return err
	}
	url := "https://" + c.baseUrl + ":" + strconv.Itoa(c.port) + shared.SendPath
	req, err := c.buildPostRequest(url, jsonData)
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

func (c *ApiClient) ReceiveMails() ([]mail.Mail, error) {
	url := "https://" + c.baseUrl + ":" + strconv.Itoa(c.port) + shared.ReceivePath
	req, err := c.buildGetRequest(url)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Returned status code: " + strconv.Itoa(resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)
	var data []mail.Mail
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (c *ApiClient) AcknowledgeMails(mailIds []string) error {
	return nil
}

func (c *ApiClient) readApiKeys(filePath string) error {
	fileExists, err := helpers.FilePathExists(filePath)
	if err != nil {
		return err
	}
	if fileExists {
		fileValue, err := helpers.ReadAllDataFromFile(filePath)
		if err != nil {
			return err
		}
		parsedData := new(persistedData)
		err = json.Unmarshal(fileValue, &parsedData)
		if err != nil {
			return err
		}
		c.receivingApiKey = parsedData.OutgoingMailApiKey
		c.sendingApiKey = parsedData.IncomingMailApiKey

	} else {
		return errors.New("api key file does not exist")
	}
	return nil
}

func CreateClient(config configuration.Configuration) (ApiClient, error) {
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
		Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: caCertPool}}}

	err = apiClient.readApiKeys(config.ApiKeysFilePath)
	if err != nil {
		return apiClient, err
	}

	return apiClient, nil
}
