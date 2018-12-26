package user

/*
	Represents a user of the ticket system.
*/
type User struct {
	Mail      string
	UserId    int
	FirstName string
	LastName  string
	Role      UserRole
}

/*
	Get the user name string.
*/
func (u *User) GetUserNameString() string {
	return u.FirstName + " " + u.LastName
}

/*
	Copy the user.
*/
func (u *User) Copy() User {
	return User{Mail: u.Mail, UserId: u.UserId, FirstName: u.FirstName, LastName: u.LastName, Role: u.Role}
}
