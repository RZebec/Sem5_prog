package user

/*
	Represents a user of the ticket system.
*/
type User struct {
	Mail      string
	UserId    int
	FirstName string
	LastName  string
}

func (u *User) GetUserNameString() string {
	return u.FirstName + " " + u.LastName
}
