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

var senders = []string{"test1@gmx.de", "Oberheld.asdf@web.de", "horstChristianAnderson@zip.org", "css@asd.com", "ad.du@ff.de",
	"kaka@lang.de", "abc.defg@hij.de", "blabla@bla.de", "labor@saft.de", "orange@blau.de", "hohlfruch@haze.de"}

var recievers = []string{"test2@gmx.de", "SuperOberheld.asdf@web.de", "FranzhorstChristianAnderson@zip.org", "Terrorist@asd.com",
	"ad.du@ff.de", "kakhaufen@lang.de", "abcbnm.defg@hij.de", "blabla@bla.de", "laborGrube@saft.de", "orange@blau.de", "hohlfruchtigerSaft@haze.de"}

func (mail Mail) RandomMail(n int) []Mail {
	mails := make([]Mail, n)
	for i := 0; i < n; i++ {
		mail.Subject = randomText(10)
		mail.Content = randomText(50)
		mail.Sender = senders[i]
		mail.Receiver = recievers[i]
		mails[i] = mail
	}
	return mails
}

func (mail Mail) ExplicitMail() []Mail {
	fmt.Print("Entry subject: ")
	mail.Subject = readEntry()
	fmt.Print("Entry text: ")
	mail.Content = readEntry()
	fmt.Print("Enter your Reciever")
	mail.Receiver = readEntry()
	fmt.Print("Enter your SenderMail")
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
