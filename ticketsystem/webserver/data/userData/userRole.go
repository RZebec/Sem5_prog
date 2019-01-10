// 5894619, 6720876, 9793350
package userData

type UserRole int

const (
	Admin UserRole = 1 + iota
	RegisteredUser
)
