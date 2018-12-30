package wrappers

import "net/http"

/*
	A test function.
*/
func (m *MockedNextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.hasBeenCalled = true
	m.request = r
	w.Write([]byte("Next handler has been called"))
	w.WriteHeader(200)
}
