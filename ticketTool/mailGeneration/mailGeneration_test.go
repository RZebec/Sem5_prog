package mailGeneration

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

/*
	The given number of mails should be generated.
*/
func TestMailGenerator_RandomMail(t *testing.T) {
	mockedIo := new(inputOutput.MockedInputOutput)
	mockedIo.On("Print", mock.Anything)
	mailGenerator := MailGenerator{io: mockedIo}
	numberOfMails := 10
	mails := mailGenerator.RandomMail(numberOfMails, 10, 50)
	assert.True(t, len(mails) == numberOfMails, "Generator do not Generate the number of mails")
}

/*
	The generated text should have the given length.
*/
func TestRandomText(t *testing.T) {
	numberOfChars := 19
	assert.Equal(t, numberOfChars, len(randomText(numberOfChars)))
}

/*
	Random generation of mail addresses should work.
*/
func TestGenerateMailAddresses(t *testing.T) {
	a, b := generateTwoMailAdresses_FromRandomPool()
	assert.True(t, a != b, "Adresses should be not the same")

	/*
		if you generate Adresses twice with same Method, the sequence of generated Emails are not the same
		=>real 'random'
	*/
	for i := 0; i < 20; i++ {
		a, b := generateTwoMailAdresses_FromRandomPool()
		fmt.Println("A: " + a + " B:" + b)
	}
	fmt.Println("__________________")
	for i := 0; i < 20; i++ {
		a, b := generateTwoMailAdresses_FromRandomPool()
		fmt.Println("A: " + a + " B:" + b)
	}
}

/*
	A explicit mail should be generated.
*/
func TestMailGenerator_ExplicitMail(t *testing.T) {
	mockedIO := new(inputOutput.MockedInputOutput)

	mockedIO.On("Print", "Entry subject: ").Once()
	mockedIO.On("ReadEntry").Return("testSubject1")
	mockedIO.On("Print", "Entry text: ").Once()
	mockedIO.On("ReadEntry").Return("testContent1")
	mockedIO.On("Print", "Enter your SenderMail: ")
	mockedIO.On("ReadEntry").Return(mock.Anything)

	testee := CreateMailGenerator(mockedIO)
	testee.ExplicitMail()

	mockedIO.AssertExpectations(t)
}
