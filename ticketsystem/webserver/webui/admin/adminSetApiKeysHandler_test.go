package admin

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

/*
	Only POST request should be possible.
*/
func TestAdminSetApiKeysHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/set_api_keys", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	testee := SetApiKeysHandler{ChangeOutgoingMailApiKey: mockedApiContext.ChangeOutgoingMailApiKey, ChangeIncomingMailApiKey: mockedApiContext.ChangeIncomingMailApiKey,
		Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Status code 405 should be returned")
	mockedApiContext.AssertExpectations(t)
}

/*
	A request with incorrect data should return a 400.
*/
func TestAdminSetApiKeysHandler_ServeHTTP_IncorrectData(t *testing.T) {
	req, err := http.NewRequest("POST", "/set_api_keys", nil)
	req.Form = url.Values{}
	req.Form.Add("incomingMailApiKey", "1234")
	req.Form.Add("outgoingMailApiKey", "4321")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)

	testee := SetApiKeysHandler{ChangeOutgoingMailApiKey: mockedApiContext.ChangeOutgoingMailApiKey, ChangeIncomingMailApiKey: mockedApiContext.ChangeIncomingMailApiKey,
		Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code 400 should be returned")
	mockedApiContext.AssertExpectations(t)
}

/*
	A valid request should be possible.
*/
func TestAdminSetApiKeysHandler_ServeHTTP_ValidRequest(t *testing.T) {
	testIncomingMailApiKey := "b0k2xW60gf3U6C5SvvnYSzTs18bf76VSkH7WSIMAwjXXF2arz7EhwTq5cSnaHN0nni4bzcoY3UW6eONFSYdHBuRSkHh1IvxPIyyrVLcZZzTAYD7SQTiWdBEVSQBu517km1"
	testOutgoingMailApiKey := "L3C2HLzf4EHae2WLezlpLL2nZVSR3LCq2H35oQQnny7MhSvHhUBYCpV3t0jTF71X6RCJ605Nv1CyQ8gTwmSQDeF11MXyjjgindFCFSC3uttoSPCR81mmj4smAtVgECThbp"

	req, err := http.NewRequest("POST", "/set_api_keys", nil)
	req.Form = url.Values{}
	req.Form.Add("incomingMailApiKey", testIncomingMailApiKey)
	req.Form.Add("outgoingMailApiKey", testOutgoingMailApiKey)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)
	mockedApiContext.On("ChangeIncomingMailApiKey", mock.Anything).Return(nil)
	mockedApiContext.On("ChangeOutgoingMailApiKey", mock.Anything).Return(nil)

	testee := SetApiKeysHandler{ChangeOutgoingMailApiKey: mockedApiContext.ChangeOutgoingMailApiKey, ChangeIncomingMailApiKey: mockedApiContext.ChangeIncomingMailApiKey,
		Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusFound, rr.Code, "Status code 302 should be returned")

	mockedApiContext.AssertExpectations(t)
}

/*
	A error from changing the outgoing api key should result in a 500.
*/
func TestAdminSetApiKeysHandler_ServeHTTP_ChangeOutgoing_ContextReturnError_500Returned(t *testing.T) {
	testIncomingMailApiKey := "b0k2xW60gf3U6C5SvvnYSzTs18bf76VSkH7WSIMAwjXXF2arz7EhwTq5cSnaHN0nni4bzcoY3UW6eONFSYdHBuRSkHh1IvxPIyyrVLcZZzTAYD7SQTiWdBEVSQBu517km1"
	testOutgoingMailApiKey := "L3C2HLzf4EHae2WLezlpLL2nZVSR3LCq2H35oQQnny7MhSvHhUBYCpV3t0jTF71X6RCJ605Nv1CyQ8gTwmSQDeF11MXyjjgindFCFSC3uttoSPCR81mmj4smAtVgECThbp"

	req, err := http.NewRequest("POST", "/set_api_keys", nil)
	req.Form = url.Values{}
	req.Form.Add("incomingMailApiKey", testIncomingMailApiKey)
	req.Form.Add("outgoingMailApiKey", testOutgoingMailApiKey)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)
	mockedApiContext.On("ChangeIncomingMailApiKey", mock.Anything).Return(nil)
	mockedApiContext.On("ChangeOutgoingMailApiKey", mock.Anything).Return(errors.New("TestError"))

	testee := SetApiKeysHandler{ChangeOutgoingMailApiKey: mockedApiContext.ChangeOutgoingMailApiKey, ChangeIncomingMailApiKey: mockedApiContext.ChangeIncomingMailApiKey,
		Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")

	mockedApiContext.AssertExpectations(t)
}

/*
	A error from changing the incoming api key should result in a 500.
*/
func TestAdminSetApiKeysHandler_ServeHTTP_ChangeIncoming_ContextReturnError_500Returned(t *testing.T) {
	testIncomingMailApiKey := "b0k2xW60gf3U6C5SvvnYSzTs18bf76VSkH7WSIMAwjXXF2arz7EhwTq5cSnaHN0nni4bzcoY3UW6eONFSYdHBuRSkHh1IvxPIyyrVLcZZzTAYD7SQTiWdBEVSQBu517km1"
	testOutgoingMailApiKey := "L3C2HLzf4EHae2WLezlpLL2nZVSR3LCq2H35oQQnny7MhSvHhUBYCpV3t0jTF71X6RCJ605Nv1CyQ8gTwmSQDeF11MXyjjgindFCFSC3uttoSPCR81mmj4smAtVgECThbp"

	req, err := http.NewRequest("POST", "/set_api_keys", nil)
	req.Form = url.Values{}
	req.Form.Add("incomingMailApiKey", testIncomingMailApiKey)
	req.Form.Add("outgoingMailApiKey", testOutgoingMailApiKey)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedApiContext := new(mockedForTests.MockedApiConfiguration)
	mockedApiContext.On("ChangeIncomingMailApiKey", mock.Anything).Return(errors.New("TestError"))

	testee := SetApiKeysHandler{ChangeOutgoingMailApiKey: mockedApiContext.ChangeOutgoingMailApiKey, ChangeIncomingMailApiKey: mockedApiContext.ChangeIncomingMailApiKey,
		Logger: testLogger}

	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")

	mockedApiContext.AssertExpectations(t)
}
