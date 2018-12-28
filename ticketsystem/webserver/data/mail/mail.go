package mail

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Mail struct {
	Id       string
	Sender   string
	Receiver string
	Subject  string
	Content  string
	SentTime time.Time
}

func (mail Mail) RandomMail(n int) []Mail {
	mails := make([]Mail, n)
	for i := 0; i < n; i++ {
		mail.Subject = randomText(10)
		mail.Content = randomText(50)
		mails[i] = mail
	}
	return mails
}

func (mail Mail) ExplicitMail() []Mail {
	fmt.Print("Entry subject: ")
	mail.Subject = readEntry()
	fmt.Print("Entry text: ")
	mail.Content = readEntry()
	fmt.Print("Enter your name")
	mail.Sender = readEntry()

	mails := make([]Mail, 1)
	mails[0] = mail
	return mails
}

func randomText(numberOfChar int) string {
	text := ""
	for i := 0; i < numberOfChar; i++ {
		r := rand.Intn(128)
		text = text + strconv.Itoa(r)
	}
	return text
}

func readEntry() string {
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	value = strings.TrimRight(value, "\n")
	return value
}
