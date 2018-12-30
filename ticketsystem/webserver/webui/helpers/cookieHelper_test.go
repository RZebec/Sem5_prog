package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetCookie(t *testing.T) {
	// Create a new HTTP Recorder (implements http.ResponseWriter)
	recorder := httptest.NewRecorder()

	// Drop a cookie on the recorder.
	SetCookie(recorder, "TestCookie", "1234")

	// Copy the Cookie over to a new Request
	request := &http.Request{Header: http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]}}

	// Extract the dropped cookie from the request.
	cookie, err := request.Cookie("TestCookie")

	if err != nil {
		t.Errorf("handler returned unexpected error: got %v", err.Error())
	}

	// Decode the cookie
	data := cookie.Value

	if data != "1234" {
		t.Errorf("handler returned unexpected cookie: got %v want %v",
			data, "1234")
	}
}