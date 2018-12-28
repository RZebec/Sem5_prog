package client

import (
	"crypto/tls"
	"crypto/x509"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
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
func CreateTestServer(t *testing.T, handler http.Handler) *httptest.Server{
	server := httptest.NewUnstartedServer(handler)
	cert, err := tls.LoadX509KeyPair("cert.pem","key.pem")
	assert.Nil(t, err)
	caCert, err := ioutil.ReadFile("cert.pem")
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
func AdjustConfigurationToTestServer(t *testing.T, conf configuration.Configuration, server httptest.Server) configuration.Configuration{
	cleanedUrl := strings.Trim(server.URL, "https://")
	cleanedUrl = strings.Trim(cleanedUrl, "http://")
	s := strings.Split(cleanedUrl, ":")
	ip, port := s[0], s[1]
	conf.BaseUrl = ip
	conf.Port, _ = strconv.Atoi(port)
	conf.CertificatePath = "cert.pem"
	return conf
}

/*
	Generates some test mails.
 */
func getTestEmails() []mail.Mail{
	var eMails []mail.Mail
	eMails = append(eMails, mail.Mail{Id: "testId1", Sender:"test1@test.de", Receiver: "testReceiver1@test.de",
		Subject: "TestSubject1", Content: "testContent1", SentTime: time.Now().Unix()})
	eMails = append(eMails, mail.Mail{Id: "testId2", Sender:"test2@test.de", Receiver: "testReceiver2@test.de",
		Subject: "TestSubject2", Content: "testContent2", SentTime: time.Now().Unix() })
	eMails = append(eMails, mail.Mail{Id: "testId3", Sender:"test3@test.de", Receiver: "testReceiver3@test.de",
		Subject: "TestSubject3", Content: "testContent3", SentTime: time.Now().Unix() })
	return eMails
}



func TestApiClient_SendMails_ReturnsOk(t *testing.T) {
	testMails := getTestEmails()

	// The handler function will assert the payload and is passed to the testserver:
	server := CreateTestServer(t, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Decode the mails:
		decoder := json.NewDecoder(req.Body)
		var mails []mail.Mail
		err := decoder.Decode(&mails)
		if err != nil {
			assert.Nil(t, err, "The received mails should be readable")
		} else {
			expectedMails := testMails
			assert.Equal(t, len(expectedMails), len(mails))
			for idx, expectedMail := range expectedMails{
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
	conf = AdjustConfigurationToTestServer(t, conf, *server)
	testee, err := CreateClient(conf)
	assert.Nil(t, err)

	// Send the mails. The handler function will assert the payload:
	err = testee.SendMails(testMails)
	assert.Nil(t, err)
}