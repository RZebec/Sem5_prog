package ticket

import "time"

/*
	A entry in a ticket.
*/
type MessageEntry struct {
	Id           int
	CreatorMail  string
	CreationTime time.Time
	Content      string
	OnlyInternal bool
}
