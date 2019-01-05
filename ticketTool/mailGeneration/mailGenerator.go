package mailGeneration

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"math/rand"
	"strconv"
	"time"
)

/*
mails for Random Mail generator
*/
var senders = []string{"test1@gmx.de", "Oberheld.asdf@web.de", "horstChristianAnderson@zip.org", "css@asd.com", "ad.du@ff.de",
	"kaka@lang.de", "abc.defg@hij.de", "blabla@bla.de", "labor@saft.de", "orange@blau.de", "hohlfruch@haze.de", "test2@gmx.de",
	"SuperOberheld.asdf@web.de", "FranzhorstChristianAnderson@zip.org", "Terrorist@asd.com",
	"ad.du@ff.de", "kakhaufen@lang.de", "abcbnm.defg@hij.de", "blabla@bla.de", "laborGrube@saft.de",
	"orange@blau.de", "hohlfruchtigerSaft@haze.de"}

type MailGeneration interface {
	RandomMail(n int, subjectLength int, contentLength int) []mail.Mail
	ExplicitMail() []mail.Mail
}

type MailGenerator struct {
	io inputOutput.InputOutput
}

/*
create Generator
*/
func CreateMailGenerator(io inputOutput.InputOutput) MailGenerator {
	return MailGenerator{io: io}
}

/*
create Random Mails by number Of Mails n, subjectLength and contentlength
*/
func (m *MailGenerator) RandomMail(n int, subjectLength int, contentLength int) []mail.Mail {
	mails := make([]mail.Mail, n)
	for i := 0; i < n; i++ {
		generatedMail := mail.Mail{}
		generatedMail.Subject = randomText(subjectLength)
		m.io.Print("Subject " + strconv.Itoa(i) + ": " + generatedMail.Subject)
		generatedMail.Content = randomText(contentLength)
		generatedMail.Sender, generatedMail.Receiver = generateTwoMailAdresses_FromRandomPool()
		generatedMail.SentTime = time.Now().Unix()
		mails[i] = generatedMail
	}
	return mails
}

//get send and recieve mailadress back
func generateTwoMailAdresses_FromRandomPool() (string, string) {
	for true {
		adressOne := senders[rand.Intn(len(senders))]
		adressTwo := senders[rand.Intn(len(senders))]
		if adressOne != adressTwo {
			return adressOne, adressTwo
		}

	}
	return "", ""
}

/*
create a Mail on your own and get back a List with one entry
*/
func (m *MailGenerator) ExplicitMail() []mail.Mail {
	email := mail.Mail{}
	m.io.Print("Entry subject: ")
	email.Subject = m.io.ReadEntry()
	m.io.Print("Entry text: ")
	email.Content = m.io.ReadEntry()
	email.Receiver = "notification@ticketsystem.de"
	m.io.Print("Enter your SenderMail: ")
	email.Sender = m.io.ReadEntry()
	email.SentTime = time.Now().Unix()

	mails := make([]mail.Mail, 1)
	mails[0] = email
	return mails
}

/*
generate Random Text in Asciichars
*/
func randomText(numberOfChar int) string {
	text := ""
	for i := 0; i < numberOfChar; i++ {
		r := rand.Intn(128)
		text = text + string(r)
	}
	return text
}
