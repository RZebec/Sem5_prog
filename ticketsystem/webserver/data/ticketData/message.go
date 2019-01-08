// 5894619, 6720876, 9793350
package ticketData

import "time"

/*
	A entry in a ticket
*/
type MessageEntry struct {
	Id           int
	CreatorMail  string
	CreationTime time.Time
	Content      string
	OnlyInternal bool
}

/*
	Copy a MessageEntry.
*/
func (s *MessageEntry) Copy() MessageEntry {
	return MessageEntry{Id: s.Id, CreatorMail: s.CreatorMail, CreationTime: s.CreationTime,
		Content: s.Content, OnlyInternal: s.OnlyInternal}
}
