package ticket

import "de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"

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
func ConvertToCreator(user user.User) Creator {
	return Creator{Mail: user.Mail, FirstName: user.FirstName, LastName: user.LastName}
}