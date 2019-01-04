package logging

import "errors"

/*
	Example for a debug message.
*/
func ExampleConsoleLogger_LogDebug() {
	logger := ConsoleLogger{}
	logger.LogDebug("My component", "My debug message")
	// Output:
	// <Debug>[My component]: My debug message
}

/*
	Example for a info message.
*/
func ExampleConsoleLogger_LogInfo() {
	logger := ConsoleLogger{}
	logger.LogInfo("My component", "My info message")
	// Output:
	// <Info>[My component]: My info message
}

/*
	Example for a error message.
*/
func ExampleConsoleLogger_LogError() {
	logger := ConsoleLogger{}
	err := errors.New("my error message")
	logger.LogError("My component", err)
	// Output:
	// <Error>[My component]: my error message
}

/*
	Example for a warning message.
*/
func ExampleConsoleLogger_LogWarning() {
	logger := ConsoleLogger{}
	logger.LogWarning("My component", "My warning message")
	// Output:
	// <Warning>[My component]: My warning message
}
