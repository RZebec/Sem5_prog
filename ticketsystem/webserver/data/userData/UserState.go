package userData

type UserState int

const (
	Active UserState = 1 + iota
	WaitingToBeUnlocked
	OnVacation
)
