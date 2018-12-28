package mail

import "time"

type Mail struct {
	Sender string
	Receiver string
	Content string
	SentTime time.Time
}