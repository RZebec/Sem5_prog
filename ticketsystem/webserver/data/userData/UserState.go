// 5894619, 6720876, 9793350
package userData

type UserState int

const (
	Active UserState = 1 + iota
	WaitingToBeUnlocked
	OnVacation
)
