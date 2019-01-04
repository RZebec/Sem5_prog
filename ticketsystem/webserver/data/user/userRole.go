package user

type UserRole int

const (
	Admin UserRole = 1 + iota
	RegisteredUser
)
