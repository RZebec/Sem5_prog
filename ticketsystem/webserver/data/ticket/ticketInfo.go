package ticket

import (
	"time"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
)

/*
	Containing only general information about a ticket.
*/
type TicketInfo struct {
	Id                   int
	Title                string
	Editor               user.User
	HasEditor			 bool
	Creator              Creator
	CreationTime         time.Time
	LastModificationTime time.Time
}
