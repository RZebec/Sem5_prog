package mail

/*
	Build the content of a notification.
*/
func BuildNotificationMailContent(receiver string, sender string, content string) string {
	prefix := "Hello " + receiver + ",\n"
	prefix = prefix + sender + " added a message to your ticket.\n"
	prefix = prefix + "Content of the message: \n"
	return prefix + content
}
