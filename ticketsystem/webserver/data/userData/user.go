// 5894619, 6720876, 9793350
package userData

/*
	Represents a user of the ticket system.
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
	Get the user name string.
*/
func (u *User) GetUserNameString() string {
	return u.FirstName + " " + u.LastName
}

/*
	Copy the user.
*/
func (u *User) Copy() User {
	return User{Mail: u.Mail, UserId: u.UserId, FirstName: u.FirstName, LastName: u.LastName, Role: u.Role, State: u.State}
}

/*
	Get a invalid default user.
*/
func GetInvalidDefaultUser() User {
	return User{Mail: "", UserId: 0, FirstName: "", LastName: "", Role: 0, State: 0}
}
