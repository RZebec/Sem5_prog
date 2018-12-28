package api

import (
	"de/vorlesung/projekt/IIIDDD/shared"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const validIncomingApiKey = "nNIr6vgamoa06F15jlnB98GGT5YY5qk4fvSsB3V8uJD3mpZAKlCMQVeEBf5SHsoXMKUBWfqKYpoj991fB3amzmmes0JaXjiTQRERXEDJsZoinD3bngqz7YjIXdNc6kll"

const validOutgoingApiKey = "zMLky9tCxQ6otKrmB3hyq2q4qnzSntW4hAVziRBuLZBh8aHJ5R7Sut72NPDGfazWDidJ0RewjYWKwKCCaBVSWCSMdafA7BWVOKFO5gBvfEj4VfIPO7cBCC0MiCbq0ZLT"

func getValidIncomingApiKey() string {
	return validIncomingApiKey
}

type MockedNextHandler struct {
	hasBeenCalled bool
	request *http.Request
}

func (m *MockedNextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	m.hasBeenCalled = true
	m.request = r
	w.WriteHeader(200)
}



func TestApiKeyAuthenticationHandler_ServeHTTP_InvalidApiKey(t *testing.T) {
	req, err := http.NewRequest("GET", shared.SendPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	nextHandler  := MockedNextHandler{}

	testee := ApiKeyAuthenticationHandler{ApiKeyResolver: getValidIncomingApiKey, Next: &nextHandler}
	handler := http.HandlerFunc(testee.ServeHTTP)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code, "Status code 401 should be returned")
	assert.False(t, nextHandler.hasBeenCalled, "Child handler should not be called")

}