package userData

/*
	Represents a userData of the ticketData system.
*/
type User struct {
	Mail      string
	UserId    int
	FirstName string
	LastName  string
	Role      UserRole
	State     UserState
}

/*
	Get the userData name string.
*/
func (u *User) GetUserNameString() string {
	return u.FirstName + " " + u.LastName
}

/*
	Copy the userData.
*/
func (u *User) Copy() User {
	return User{Mail: u.Mail, UserId: u.UserId, FirstName: u.FirstName, LastName: u.LastName, Role: u.Role, State: u.State}
}

/*
	Get a invalid default userData.
*/
func GetInvalidDefaultUser() User {
	return User{Mail: "", UserId: 0, FirstName: "", LastName: "", Role: 0, State: 0}
}
