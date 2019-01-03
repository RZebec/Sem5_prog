package testhelpers

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"net/http"
)

/*
	Get a test logger.
 */
func GetTestLogger() logging.Logger {
	return logging.ConsoleLogger{SetTimeStamp: false}
}


/*
	A http handler which stores the request.
 */
type LoggingHTPPHandler struct {
	HasBeenCalled bool
	Request *http.Request
}

/*
	The method to capture the request.
 */
func (h *LoggingHTPPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.HasBeenCalled = true
	h.Request = r
	w.WriteHeader(200)
}