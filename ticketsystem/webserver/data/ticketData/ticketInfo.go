// 5894619, 6720876, 9793350
package ticketData

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"time"
)

/*
	Containing only general information about a ticket.
*/
type TicketInfo struct {
	Id                   int
	Title                string
	Editor               userData.User
	HasEditor            bool
	Creator              Creator
	CreationTime         time.Time
	LastModificationTime time.Time
	State                TicketState
}

/*
	Copy the TicketInfo.
*/
func (s *TicketInfo) Copy() TicketInfo {
	return TicketInfo{Id: s.Id, Title: s.Title, Editor: s.Editor.Copy(), HasEditor: s.HasEditor,
		Creator: s.Creator.Copy(), CreationTime: s.CreationTime,
		LastModificationTime: s.LastModificationTime, State: s.State}
}
