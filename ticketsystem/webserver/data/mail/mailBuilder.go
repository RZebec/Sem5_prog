package mail

import "strconv"

/*
	Build the content of a notification.
*/
func BuildAppendMessageNotificationMailContent(receiver string, sender string, content string) string {
	prefix := "Hello " + receiver + ",\n"
	prefix = prefix + sender + " added a message to your ticket.\n"
	prefix = prefix + "Content of the message: \n"
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
	prefix := "Hello " + receiver + ",\n"
	prefix = prefix + "Your account has been activated by the administrator."
	return prefix
}

/*
	Build the subject string.
 */
func BuildAppendMessageNotificationMailSubject(ticketId int) string {
	stringValue := strconv.Itoa(ticketId)
	return "A new message has been appended to your ticket with TicketId<" + stringValue + ">:"
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
	prefix := "Hello " + receiver + ",\n"
	prefix = prefix + "Ticket " + firstTicket + " has been merged with " + secondTicket + "."
	return prefix
}
