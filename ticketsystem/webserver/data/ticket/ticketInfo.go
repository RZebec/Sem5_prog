package ticket

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"time"
)

/*
	Containing only general information about a ticket.
*/
type TicketInfo struct {
	Id                   int
	Title                string
	Editor               user.User
	HasEditor            bool
	Creator              Creator
	CreationTime         time.Time
	LastModificationTime time.Time
}

/*
	Copy the TicketInfo.
 */
func (s *TicketInfo) Copy() (TicketInfo){
	return TicketInfo{Id: s.Id, Title: s.Title, Editor: s.Editor.Copy(), HasEditor: s.HasEditor,
		Creator: s.Creator.Copy(), CreationTime: s.CreationTime,
		LastModificationTime: s.LastModificationTime}
}
