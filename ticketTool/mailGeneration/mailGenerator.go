package mailGeneration

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/validation/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"math/rand"
	"strconv"
	"time"
)

/*
mails for Random Mail generator
*/
var senders = []string{"test1@gmx.de", "Oberheld.asdf@web.de", "horstChristianAnderson@zip.org", "css@asd.com", "ad.du@ff.de",
	"test@lang.de", "abc.defg@hij.de", "blabla@bla.de", "labor@saft.de", "orange@blau.de", "hohlfruch@haze.de", "test2@gmx.de",
	"SuperOberheld.asdf@web.de", "FranzhorstChristianAnderson@zip.org", "Terrorist@asd.com",
	"ad.du@ff.de", "Müller@lang.de", "abcbnm.defg@hij.de", "blabla@bla.de", "laborGrube@saft.de",
	"orange@blau.de", "hohlfruchtigerSaft@haze.de"}

/*
	A interface for the mail generation.
*/
type MailGeneration interface {
	RandomMail(n int, subjectLength int, contentLength int) []mailData.Mail
	ExplicitMail() []mailData.Mail
}

/*
	A mail generator.
*/
type MailGenerator struct {
	io inputOutput.InputOutput
}

/*
	Create Generator
*/
func CreateMailGenerator(io inputOutput.InputOutput) MailGenerator {
	return MailGenerator{io: io}
}

/*
	Create Random Mails by number Of Mails n, subjectLength and contentlength
*/
func (m *MailGenerator) RandomMail(n int, subjectLength int, contentLength int) []mailData.Mail {
	mails := make([]mailData.Mail, n)
	for i := 0; i < n; i++ {
		generatedMail := mailData.Mail{}
		generatedMail.Subject = randomText(subjectLength)
		m.io.Print("Subject " + strconv.Itoa(i) + ": " + generatedMail.Subject)
		generatedMail.Content = randomText(contentLength)
		generatedMail.Sender, generatedMail.Receiver = generateTwoMailAdresses_FromRandomPool()
		generatedMail.SentTime = time.Now().Unix()
		mails[i] = generatedMail
	}
	return mails
}

/*
	Get two random mail addresses.
*/
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
func (m *MailGenerator) ExplicitMail() []mailData.Mail {
	validator := mail.NewValidator()
	email := mailData.Mail{}
	m.io.Print("Entry subject: ")
	email.Subject = m.io.ReadEntry()
	m.io.Print("Entry text: ")
	email.Content = m.io.ReadEntry()
	email.Receiver = "notification@ticketsystem.de"
	for true {
		m.io.Print("Enter your Sender-Mailadress: ")
		sendAdress := m.io.ReadEntry()
		if validator.Validate(sendAdress) {
			email.Sender = sendAdress
			break
		} else {
			m.io.Print("This is not a valide Adress. Retry!")
		}
	}

	email.SentTime = time.Now().Unix()

	mails := make([]mailData.Mail, 1)
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
