package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAdminSetApiKeysHandlerWrongRequestMethod_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("GET", "/set_api_keys", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := getTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("ChangeIncomingMailApiKey", mock.Anything)

	mockedApiContext.On("ChangeOutgoingMailApiKey", mock.Anything)

	testee := AdminSetApiKeysHandler{ApiConfiguration: mockedApiContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Status code 405 should be returned")
}

func TestAdminSetApiKeysHandlerIncorrectData_ServeHTTP(t *testing.T) {
	data := url.Values{}
	data.Set("incomingMailApiKey", "1234")
	data.Set("outgoingMailApiKey", "4321")

	req, err := http.NewRequest("POST", "/set_api_keys",strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := getTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("ChangeIncomingMailApiKey", mock.Anything)

	mockedApiContext.On("ChangeOutgoingMailApiKey", mock.Anything)

	testee := AdminSetApiKeysHandler{ApiConfiguration: mockedApiContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code 400 should be returned")
}

func TestAdminSetApiKeysHandler_ServeHTTP(t *testing.T) {
	testIncomingMailApiKey := "b0k2xW60gf3U6C5SvvnYSzTs18bf76VSkH7WSIMAwjXXF2arz7EhwTq5cSnaHN0nni4bzcoY3UW6eONFSYdHBuRSkHh1IvxPIyyrVLcZZzTAYD7SQTiWdBEVSQBu517km1"
	testOutgoingMailApiKey := "L3C2HLzf4EHae2WLezlpLL2nZVSR3LCq2H35oQQnny7MhSvHhUBYCpV3t0jTF71X6RCJ605Nv1CyQ8gTwmSQDeF11MXyjjgindFCFSC3uttoSPCR81mmj4smAtVgECThbp"

	data := url.Values{}
	data.Add("incomingMailApiKey", testIncomingMailApiKey)
	data.Add("outgoingMailApiKey", testOutgoingMailApiKey)

	req, err := http.NewRequest("POST", "/set_api_keys",strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := getTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	mockedApiContext.On("ChangeIncomingMailApiKey", mock.Anything)

	mockedApiContext.On("ChangeOutgoingMailApiKey", mock.Anything)

	testee := AdminSetApiKeysHandler{ApiConfiguration: mockedApiContext, Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200ä should be returned")
}
