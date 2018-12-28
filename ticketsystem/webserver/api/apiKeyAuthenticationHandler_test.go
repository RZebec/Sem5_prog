package api

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const validApiKey = "nNIr6vgamoa06F15jlnB98GGT5YY5qk4fvSsB3V8uJD3mpZAKlCMQVeEBf5SHsoXMKUBWfqKYpoj991fB3amzmmes0JaXjiTQRERXEDJsZoinD3bngqz7YjIXdNc6kll"

func getValidIncomingApiKey() string {
	return validApiKey
}

// A mocked handler for tests
type MockedNextHandler struct {
	hasBeenCalled bool
	request       *http.Request
}

/*
	A test function.
*/
func (m *MockedNextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.hasBeenCalled = true
	m.request = r
	w.Write([]byte("Next handler has been called"))
	w.WriteHeader(200)
}

/*
	Request with no api key should return a 401.
*/
func TestApiKeyAuthenticationHandler_ServeHTTP_NoApiKey(t *testing.T) {
	req, err := http.NewRequest("POST", shared.SendPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	nextHandler := MockedNextHandler{}

	testee := ApiKeyAuthenticationHandler{ApiKeyResolver: getValidIncomingApiKey, Next: &nextHandler, AllowedMethod: "POST"}
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code, "Status code 401 should be returned")
	assert.False(t, nextHandler.hasBeenCalled, "Child handler should not be called")

}

/*
	Request with a wrong api key should return a 401.
*/
func TestApiKeyAuthenticationHandler_ServeHTTP_WrongApiKey(t *testing.T) {
	req, err := http.NewRequest("POST", shared.SendPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	req.Header.Set("Cookie", shared.AuthenticationCookieName+"=1234568")

	nextHandler := MockedNextHandler{}

	testee := ApiKeyAuthenticationHandler{ApiKeyResolver: getValidIncomingApiKey, Next: &nextHandler, AllowedMethod: "POST"}
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code, "Status code 401 should be returned")
	assert.False(t, nextHandler.hasBeenCalled, "Child handler should not be called")
}

/*
	Calling the authentication handler with the correct api key should transfer the request to the child handler.
*/
func TestApiKeyAuthenticationHandler_ServeHTTP_CorrectApiKey_NextHandlerCalled(t *testing.T) {
	req, err := http.NewRequest("POST", shared.SendPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	req.Header.Set("Cookie", shared.AuthenticationCookieName+"="+getValidIncomingApiKey())

	nextHandler := MockedNextHandler{}

	testee := ApiKeyAuthenticationHandler{ApiKeyResolver: getValidIncomingApiKey, Next: &nextHandler, AllowedMethod: "POST"}
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")
	assert.True(t, nextHandler.hasBeenCalled, "Child handler should not be called")
	assert.Equal(t, "Next handler has been called", rr.Body.String())
}


/*
	Only Post method should be accepted.
*/
func TestApiKeyAuthenticationHandler_ServeHTTP_GetRequest(t *testing.T) {
	req, err := http.NewRequest("GET", shared.SendPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	req.Header.Set("Cookie", shared.AuthenticationCookieName+"="+getValidIncomingApiKey())

	nextHandler := MockedNextHandler{}

	testee := ApiKeyAuthenticationHandler{ApiKeyResolver: getValidIncomingApiKey, Next: &nextHandler, AllowedMethod: "POST"}
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code, "Status code 401 should be returned")
	assert.False(t, nextHandler.hasBeenCalled, "Child handler should not be called")
}