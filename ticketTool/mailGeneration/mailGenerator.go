package mailGeneration

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
	"math/rand"
	"strconv"
)

var senders = []string{"test1@gmx.de", "Oberheld.asdf@web.de", "horstChristianAnderson@zip.org", "css@asd.com", "ad.du@ff.de",
	"kaka@lang.de", "abc.defg@hij.de", "blabla@bla.de", "labor@saft.de", "orange@blau.de", "hohlfruch@haze.de"}

var recievers = []string{"test2@gmx.de", "SuperOberheld.asdf@web.de", "FranzhorstChristianAnderson@zip.org", "Terrorist@asd.com",
	"ad.du@ff.de", "kakhaufen@lang.de", "abcbnm.defg@hij.de", "blabla@bla.de", "laborGrube@saft.de", "orange@blau.de", "hohlfruchtigerSaft@haze.de"}

type MailGenerator struct {
}

func (m *MailGenerator) RandomMail(n int) []mail.Mail {
	mails := make([]mail.Mail, n)
	for i := 0; i < n; i++ {
		mail := mail.Mail{}
		mail.Subject = randomText(10)
		mail.Content = randomText(50)
		mail.Sender = senders[i]
		mail.Receiver = recievers[i]
		mails[i] = mail
	}
	return mails
}

func (m *MailGenerator) ExplicitMail() []mail.Mail {
	email := mail.Mail{}
	fmt.Print("Entry subject: ")
	email.Subject = inputOutput.ReadEntry()
	fmt.Print("Entry text: ")
	email.Content = inputOutput.ReadEntry()
	fmt.Print("Enter your Reciever: ")
	email.Receiver = inputOutput.ReadEntry()
	fmt.Print("Enter your SenderMail: ")
	email.Sender = inputOutput.ReadEntry()

	mails := make([]mail.Mail, 1)
	mails[0] = email
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
