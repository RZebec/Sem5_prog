package mails

import "de/vorlesung/projekt/IIIDDD/ticketsystem/logging"

/*
	A logger for tests.
*/
func getTestLogger() logging.Logger {
	return logging.ConsoleLogger{SetTimeStamp: false}
}
