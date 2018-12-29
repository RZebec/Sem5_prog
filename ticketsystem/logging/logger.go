// Contains functions for logging purposes.
package logging

import (
	"fmt"
	"time"
)

/*
	The logger interface.
*/
type Logger interface {
	LogDebug(prefix string, message string)
	LogInfo(prefix string, message string)
	LogWarning(prefix string, message string)
	LogError(prefix string, err error)
}

/*
	A logger which sends the messages to the console.
*/
type ConsoleLogger struct {
	SetTimeStamp bool
}

/*
	Log a debug message.
*/
func (l ConsoleLogger) LogDebug(prefix string, message string) {
	timestamp := ""
	if l.SetTimeStamp {
		timestamp = time.Now().Format("2006-01-02 15:04:05") + ":"
	}
	fmt.Printf("%v<%v>[%v]: %v\n", timestamp, "Debug", prefix, message)
}

/*
	Log a info message.
*/
func (l ConsoleLogger) LogInfo(prefix string, message string) {
	timestamp := ""
	if l.SetTimeStamp {
		timestamp = time.Now().Format("2006-01-02 15:04:05") + ":"
	}
	fmt.Printf("%v<%v>[%v]: %v\n", timestamp, "Info", prefix, message)
}

/*
	Log a warning.
*/
func (l ConsoleLogger) LogWarning(prefix string, message string) {
	timestamp := ""
	if l.SetTimeStamp {
		timestamp = time.Now().Format("2006-01-02 15:04:05") + ":"
	}
	fmt.Printf("%v<%v>[%v]: %v\n", timestamp, "Warning", prefix, message)
}

/*
	Log an error.
*/
func (l ConsoleLogger) LogError(prefix string, err error) {
	if err == nil {
		return
	}
	timestamp := ""
	if l.SetTimeStamp {
		timestamp = time.Now().Format("2006-01-02 15:04:05") + ":"
	}
	fmt.Printf("%v<%v>[%v]: %v\n", timestamp, "Error", prefix, err.Error())
}
