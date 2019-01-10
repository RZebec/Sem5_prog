// 5894619, 6720876, 9793350
package sender

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/mailGeneration"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"github.com/stretchr/testify/mock"
	"testing"
)

func getTestMails() []mailData.Mail {
	var testMails []mailData.Mail
	testMails = append(testMails, mailData.Mail{Id: "testId1", Sender: "test@test.de", Receiver: "testReceiver1@test.de",
		Subject: "testSubject1", Content: "testContent1"})
	testMails = append(testMails, mailData.Mail{Id: "testId2", Sender: "test@test.de", Receiver: "testReceiver2@test.de",
		Subject: "testSubject2", Content: "testContent2"})
	testMails = append(testMails, mailData.Mail{Id: "testId3", Sender: "test@test.de", Receiver: "testReceiver3@test.de",
		Subject: "testSubject3", Content: "testContent3"})
	return testMails
}

func TestSender_SendExplicitMail(t *testing.T) {
	testMails := getTestMails()
	config := configuration.Configuration{}
	mockedIO := new(inputOutput.MockedInputOutput)
	mockedApiClient := new(client.MockedClient)
	mockedMailGenerator := new(mailGeneration.MockedMailGenerator)

	mockedIO.On("ReadEntry").Return("e")
	mockedIO.On("Print", mock.Anything)
	mockedMailGenerator.On("ExplicitMail").Return(testMails)

	mockedApiClient.On("SendMails", testMails).Return(nil)

	testee := CreateSender(config, mockedIO, mockedApiClient, mockedMailGenerator)
	testee.Run()

	mockedIO.AssertExpectations(t)
	mockedMailGenerator.AssertExpectations(t)
	mockedApiClient.AssertExpectations(t)
}

func TestSender_SendRandomMails(t *testing.T) {

	testMails := getTestMails()
	config := configuration.Configuration{}
	mockedIO := new(inputOutput.MockedInputOutput)
	mockedApiClient := new(client.MockedClient)
	mockedMailGenerator := new(mailGeneration.MockedMailGenerator)

	mockedIO.On("ReadEntry").Return("r").Once()
	mockedIO.On("ReadEntry").Return("5")
	mockedIO.On("Print", mock.Anything)
	//mockedMailGenerator.On("ExplicitMail").Return(testMails)
	mockedMailGenerator.On("RandomMail", 5, 10, 50).Return(testMails)

	mockedApiClient.On("SendMails", testMails).Return(nil)

	testee := CreateSender(config, mockedIO, mockedApiClient, mockedMailGenerator)
	testee.Run()

	mockedIO.AssertExpectations(t)
	mockedMailGenerator.AssertExpectations(t)
	mockedApiClient.AssertExpectations(t)
}
