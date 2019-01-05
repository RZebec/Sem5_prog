package mailData

/*
	A email.
*/
type Mail struct {
	Id       string
	Sender   string
	Receiver string
	Subject  string
	Content  string
	SentTime int64
	Headers  []string
}
