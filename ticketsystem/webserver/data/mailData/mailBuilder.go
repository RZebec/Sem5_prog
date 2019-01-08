// 5894619, 6720876, 9793350
package mailData

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticketData"
	"strconv"
)

/*
	Build the content of a notification.
*/
func BuildAppendMessageNotificationMailContent(receiver string, sender string, content string) string {
	prefix := "Hello " + receiver + ",\n "
	prefix = prefix + sender + " added a message to your ticket.\n "
	prefix = prefix + "Content of the message: \n "
	return prefix + content
}

/*
	Build the subject string.
*/
func BuildUnlockUserNotificationMailSubject() string {
	return "Your account has been unlocked."
}

/*
	Build the content of a notification.
*/
func BuildUnlockUserNotificationMailContent(receiver string) string {
	prefix := "Hello " + receiver + ",\n "
	prefix = prefix + "Your account has been activated by the administrator"
	return prefix
}

/*
	Build the subject string.
*/
func BuildAppendMessageNotificationMailSubject(ticketId int) string {
	stringValue := strconv.Itoa(ticketId)
	return "A new message has been appended to your ticket with TicketId<" + stringValue + ">"
}

/*
	Build the subject string.
*/
func BuildTicketMergeNotificationMailSubject(ticketId int, newTicketId int) string {
	stringValue := strconv.Itoa(ticketId)
	newTicketIdValue := strconv.Itoa(newTicketId)
	return "Your ticket with TicketId<" + stringValue + "> has been merged. New Ticket: TicketId<" + newTicketIdValue + ">"
}

/*
	Build the content of a notification.
*/
func BuildTicketMergeNotificationMailContent(receiver string, firstTicketId int, secondTicketId int) string {
	firstTicket := strconv.Itoa(firstTicketId)
	secondTicket := strconv.Itoa(secondTicketId)
	prefix := "Hello " + receiver + ",\n "
	prefix = prefix + "Ticket " + firstTicket + " has been merged with " + secondTicket + "."
	return prefix
}

/*
	Build the subject string.
*/
func BuildTicketEditorChangedNotificationMailSubject(ticketId int) string {
	stringValue := strconv.Itoa(ticketId)
	return "Your ticket with TicketId<" + stringValue + "> has been changed. A new editor has been set or removed"
}

/*
	Build the content of a notification.
*/
func BuildTicketEditorChangedNotificationMailContent(receiver string, ticketId int, newEditor string) string {
	ticketIdString := strconv.Itoa(ticketId)
	prefix := "Hello " + receiver + ",\n "
	prefix = prefix + "Ticket " + ticketIdString + " has a new Editor: " + newEditor + "."
	return prefix
}

/*
	Build the content of a notification.
*/
func BuildTicketEditorRemovedNotificationMailContent(receiver string, ticketId int) string {
	ticketIdString := strconv.Itoa(ticketId)
	prefix := "Hello " + receiver + ",\n "
	prefix = prefix + "The editor has been removed from Ticket " + ticketIdString + "."
	return prefix
}

/*
	Build the state change notification subject.
*/
func BuildTicketStateChangedNotificationMailSubject(ticketId int) string {
	stringValue := strconv.Itoa(ticketId)
	return "The state of your ticket with TicketId<" + stringValue + "> has been changed."
}

/*
Build the state change notification content.
*/
func BuildTicketStateChangedNotificationMailContent(receiver string, state ticketData.TicketState) string {
	prefix := "Hello " + receiver + ",\n "
	prefix = prefix + "The state of your ticket has been changed to " + state.String() + "."
	return prefix
}
