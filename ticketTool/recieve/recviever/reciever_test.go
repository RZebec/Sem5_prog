package recviever

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/confirm"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"github.com/stretchr/testify/mock"
	"testing"
)

func getTestAcknowledges() []mail.Acknowledgment {
	var testAcknowledges []mail.Acknowledgment
	testAcknowledges = append(testAcknowledges, mail.Acknowledgment{Id: "testId1", Subject: "testSubject1"})
	testAcknowledges = append(testAcknowledges, mail.Acknowledgment{Id: "testId2", Subject: "testSubject2"})
	testAcknowledges = append(testAcknowledges, mail.Acknowledgment{Id: "testId3", Subject: "testSubject3"})
	return testAcknowledges
}

func getTestMails() []mail.Mail {
	var testMails []mail.Mail
	testMails = append(testMails, mail.Mail{Id: "testId1", Sender: "test@test.de", Receiver: "testReceiver1@test.de",
		Subject: "testSubject1", Content: "testContent1"})
	testMails = append(testMails, mail.Mail{Id: "testId2", Sender: "test@test.de", Receiver: "testReceiver2@test.de",
		Subject: "testSubject2", Content: "testContent2"})
	testMails = append(testMails, mail.Mail{Id: "testId3", Sender: "test@test.de", Receiver: "testReceiver3@test.de",
		Subject: "testSubject3", Content: "testContent3"})
	return testMails
}

func Test_AcknowledgeAll(t *testing.T) {
	testAcknowledges := getTestAcknowledges()
	config := configuration.Configuration{}
	mockedIO := new(inputOutput.MockedInputOutput)
	mockedApiClient := new(client.MockedClient)
	mockedConfirm := new(confirm.MockedConfirm)

	mockedIO.On("ReadEntry").Return("all")
	mockedIO.On("Print").Return(mock.Anything)
	mockedConfirm.On("GetAllAcknowledges").Return(testAcknowledges)
}

func Test_AcknowledgeSpecify(t *testing.T) {

}
