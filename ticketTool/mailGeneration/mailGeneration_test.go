package mailGeneration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMailGenerator_RandomMail(t *testing.T) {
	mailGenerator:=MailGenerator{}
	numberOfMails:=10
	mails:=mailGenerator.RandomMail(numberOfMails,10,50)
	assert.True(t,len(mails)==numberOfMails,"Generator do not Generate the number of mails")
}

func TestRandomText(t *testing.T){
	numberOfChars:=19
	assert.Equal(t,numberOfChars,len(randomText(numberOfChars)))
}

func TestGenerateMailAdresses(t *testing.T){
	a,b := generateTwoMailAdresses_FromRandomPool()
	assert.True(t,a!=b,"Adresses should be not the same")

	/*
	if you generate Adresses twice with same Method, the sequence of generated Emails are not the same
	=>real 'random'
	 */
	for i:=0;i<20;i++{
		a,b := generateTwoMailAdresses_FromRandomPool()
		fmt.Println("A: "+a+" B:"+b)
	}
	fmt.Println("__________________")
	for i:=0;i<20;i++{
		a,b := generateTwoMailAdresses_FromRandomPool()
		fmt.Println("A: "+a+" B:"+b)
	}
}
