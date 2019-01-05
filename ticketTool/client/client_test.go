package client

import (
	"crypto/tls"
	"crypto/x509"
	"de/vorlesung/projekt/IIIDDD/shared"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

/*
	Create a test server using https and the provided certificates.
*/
func CreateTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	server := httptest.NewUnstartedServer(handler)
	cert, err := tls.LoadX509KeyPair("test_cert.pem", "test_key.pem")
	assert.Nil(t, err)
	caCert, err := ioutil.ReadFile("test_cert.pem")
	assert.Nil(t, err)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}
	server.TLS = tlsConfig

	return server
}

/*
	Adjust the configuration for the client to the settings, used by the test server.
*/
func AdjustConfigurationToTestServer(t *testing.T, conf configuration.Configuration, server httptest.Server) configuration.Configuration {
	cleanedUrl := strings.Trim(server.URL, "https://")
	cleanedUrl = strings.Trim(cleanedUrl, "http://")
	s := strings.Split(cleanedUrl, ":")
	ip, port := s[0], s[1]
	conf.BaseUrl = ip
	conf.Port, _ = strconv.Atoi(port)
	conf.CertificatePath = "test_cert.pem"
	return conf
}

/*
	Generates some test mails.
*/
func getTestEmails() []mailData.Mail {
	var eMails []mailData.Mail
	eMails = append(eMails, mailData.Mail{Id: "testId1", Sender: "test1@test.de", Receiver: "testReceiver1@test.de",
		Subject: "TestSubject1", Content: "testContent1", SentTime: time.Now().Unix()})
	eMails = append(eMails, mailData.Mail{Id: "testId2", Sender: "test2@test.de", Receiver: "testReceiver2@test.de",
		Subject: "TestSubject2", Content: "testContent2", SentTime: time.Now().Unix()})
	eMails = append(eMails, mailData.Mail{Id: "testId3", Sender: "test3@test.de", Receiver: "testReceiver3@test.de",
		Subject: "TestSubject3", Content: "testContent3", SentTime: time.Now().Unix()})
	return eMails
}

/*
	Sending emails to the api should work.
*/
func TestApiClient_SendMails_ReturnsOk(t *testing.T) {
	testMails := getTestEmails()

	// The handler function will assert the payload and is passed to the testserver:
	server := CreateTestServer(t, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Decode the mails:
		decoder := json.NewDecoder(req.Body)
		var mails []mailData.Mail
		err := decoder.Decode(&mails)
		if err != nil {
			assert.Nil(t, err, "The received mails should be readable")
		} else {
			expectedMails := testMails
			assert.Equal(t, len(expectedMails), len(mails))
			for idx, expectedMail := range expectedMails {
				actualMail := mails[idx]
				assert.Equal(t, expectedMail, actualMail, "The received mails should be equal to the sent mails")
			}
		}

		rw.WriteHeader(200)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Start the server:
	server.StartTLS()

	// Configure the client:
	conf := configuration.Configuration{}
	conf.ApiKeysFilePath = "test_api.keys"
	conf = AdjustConfigurationToTestServer(t, conf, *server)
	testee, err := CreateClient(conf)
	assert.Nil(t, err)

	// Send the mails. The handler function will assert the payload:
	err = testee.SendMails(testMails)
	assert.Nil(t, err)
}

/*
	A post request should be prepared, with the payload and api keys set.
*/
func TestApiClient_preparePostRequest(t *testing.T) {
	conf := configuration.Configuration{}
	conf.CertificatePath = "test_cert.pem"
	conf.ApiKeysFilePath = "test_api.keys"
	testee, err := CreateClient(conf)
	assert.Nil(t, err)

	url := "testUrl/test"
	testpayload := []byte("myTestPayload")
	req, err := testee.buildPostRequest(url, testpayload)
	assert.Nil(t, err)
	assert.Equal(t, "POST", req.Method, "The request method should be set to POST")

	requestPayload, _ := ioutil.ReadAll(req.Body)
	assert.Equal(t, string(testpayload), string(requestPayload), "The payload should be set")

	apikeySet := false
	cookies := req.Cookies()
	for _, cookie := range cookies {
		if strings.ToLower(cookie.Name) == strings.ToLower(shared.AuthenticationCookieName) {
			apiKey := cookie.Value
			if apiKey == testee.sendingApiKey {
				apikeySet = true
				break
			}

		}
	}
	assert.True(t, apikeySet, "The api key must be set")
}

/*
	A get request should be prepared, with the api keys set.
*/
func TestApiClient_prepareGetRequest(t *testing.T) {
	conf := configuration.Configuration{}
	conf.CertificatePath = "test_cert.pem"
	conf.ApiKeysFilePath = "test_api.keys"
	testee, err := CreateClient(conf)
	assert.Nil(t, err)

	url := "testUrl/test"
	req, err := testee.buildGetRequest(url)
	assert.Nil(t, err)
	assert.Equal(t, "GET", req.Method, "The request method should be set to GET")

	assert.Equal(t, int64(0), req.ContentLength, "The payload should be empty")

	apikeySet := false
	cookies := req.Cookies()
	for _, cookie := range cookies {
		if strings.ToLower(cookie.Name) == strings.ToLower(shared.AuthenticationCookieName) {
			apiKey := cookie.Value
			if apiKey == testee.receivingApiKey {
				apikeySet = true
				break
			}

		}
	}
	assert.True(t, apikeySet, "The api key must be set")
}

/*
	Receiving emails to the api should work.
*/
func TestApiClient_ReceiveMails_MailsReturned(t *testing.T) {
	testMails := getTestEmails()

	// The handler function will return the test emails
	server := CreateTestServer(t, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		jsonData, err := json.Marshal(testMails)
		if err != nil {
			rw.WriteHeader(500)
			return
		}
		rw.Write(jsonData)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Start the server:
	server.StartTLS()

	// Configure the client:
	conf := configuration.Configuration{}
	conf.ApiKeysFilePath = "test_api.keys"
	conf = AdjustConfigurationToTestServer(t, conf, *server)
	testee, err := CreateClient(conf)
	assert.Nil(t, err)

	// Receive the mails.
	mails, err := testee.ReceiveMails()
	assert.Nil(t, err)
	for idx, expectedMail := range testMails {
		actualMail := mails[idx]
		assert.Equal(t, expectedMail, actualMail, "All mails should be transmitted and equal to the sent one")
	}
}

/*
	Sending acknowledgments to the api should work.
*/
func TestApiClient_AcknowledgeMails_ReturnsOk(t *testing.T) {
	testAcks := getTestAcknowledgements()

	// The handler function will assert the payload and is passed to the testserver:
	server := CreateTestServer(t, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Decode the mails:
		decoder := json.NewDecoder(req.Body)
		var mails []mailData.Acknowledgment
		err := decoder.Decode(&mails)
		if err != nil {
			assert.Nil(t, err, "The received acknowledgements should be readable")
		} else {
			expectedAcks := testAcks
			assert.Equal(t, len(expectedAcks), len(mails))
			for idx, expectedMail := range expectedAcks {
				actualMail := mails[idx]
				assert.Equal(t, expectedMail, actualMail,
					"The acknowledgements mails should be equal to the sent acknowledgements")
			}
		}

		rw.WriteHeader(200)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Start the server:
	server.StartTLS()

	// Configure the client:
	conf := configuration.Configuration{}
	conf.ApiKeysFilePath = "test_api.keys"
	conf = AdjustConfigurationToTestServer(t, conf, *server)
	testee, err := CreateClient(conf)
	assert.Nil(t, err)

	// Send the mails. The handler function will assert the payload:
	err = testee.AcknowledgeMails(testAcks)
	assert.Nil(t, err)
}

/*
	Get Acknowledgments for tests.
*/
func getTestAcknowledgements() []mailData.Acknowledgment {
	var acks []mailData.Acknowledgment
	acks = append(acks, mailData.Acknowledgment{Id: "testId1", Subject: "TestSubject1"})
	acks = append(acks, mailData.Acknowledgment{Id: "testId2", Subject: "TestSubject2"})
	acks = append(acks, mailData.Acknowledgment{Id: "testId3", Subject: "TestSubject3"})
	return acks
}
