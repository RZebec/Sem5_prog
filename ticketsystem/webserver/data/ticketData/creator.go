// 5894619, 6720876, 9793350
package ticketData

import "de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"

/*
	Represents the creator of a ticket.
*/
type Creator struct {
	Mail      string
	FirstName string
	LastName  string
}

/*
	Convert a user to a creator.
*/
func ConvertToCreator(user userData.User) Creator {
	return Creator{Mail: user.Mail, FirstName: user.FirstName, LastName: user.LastName}
}

/*
	Copy the creator struct.
*/
func (s *Creator) Copy() Creator {
	return Creator{Mail: s.Mail, FirstName: s.FirstName, LastName: s.LastName}
}
