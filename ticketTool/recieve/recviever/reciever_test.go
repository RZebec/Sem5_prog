package recviever

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/acknowledgementStorage"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/confirm"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"errors"
	"strconv"
	"testing"
)

func getTestAcknowledges() []mailData.Acknowledgment {
	var testAcknowledges []mailData.Acknowledgment
	testAcknowledges = append(testAcknowledges, mailData.Acknowledgment{Id: "testId1", Subject: "testSubject1"})
	testAcknowledges = append(testAcknowledges, mailData.Acknowledgment{Id: "testId2", Subject: "testSubject2"})
	testAcknowledges = append(testAcknowledges, mailData.Acknowledgment{Id: "testId3", Subject: "testSubject3"})
	return testAcknowledges
}
func getSpecifyAcknowledges() []mailData.Acknowledgment {
	var testAcknowledges []mailData.Acknowledgment
	testAcknowledges = append(testAcknowledges, mailData.Acknowledgment{Id: "testId2", Subject: "testSubject2"})
	testAcknowledges = append(testAcknowledges, mailData.Acknowledgment{Id: "testId3", Subject: "testSubject3"})
	return testAcknowledges
}

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
func getTestSelectedAcknowledge() []mailData.Acknowledgment {
	var testAcknowledges []mailData.Acknowledgment
	testAcknowledges = append(testAcknowledges, mailData.Acknowledgment{Id: "testId1", Subject: "testSubject1"})
	return testAcknowledges
}

func Test_AcknowledgeAll(t *testing.T) {
	testAcknowledges := getTestAcknowledges()
	testMails := getTestMails()
	config := configuration.Configuration{}
	mockedIO := new(inputOutput.MockedInputOutput)
	mockedApiClient := new(client.MockedClient)
	mockedConfirm := new(confirm.MockedConfirm)
	mockedStorage := new(acknowledgementStorage.MockedAcknowledgementStorage)

	mockedApiClient.On("ReceiveMails").Return(testMails, nil)
	mockedIO.On("Print", strconv.Itoa(len(testMails))+" Mails are coming from Server")
	mockedConfirm.On("GetAllAcknowledges", testMails).Return(testAcknowledges)
	mockedStorage.On("AppendAcknowledgements", testAcknowledges).Return(nil)
	mockedIO.On("Print", "Save Acknowledges...")
	mockedStorage.On("ReadAcknowledgements").Return(testAcknowledges, nil)
	mockedIO.On("Print", "Available Mails: "+strconv.Itoa(len(testAcknowledges)))
	mockedIO.On("Print", "send all Acknowledges or specify Acknowledges to Server. Or stop reciever (all/specify/stop):")
	mockedIO.On("ReadEntry").Return("all")
	mockedApiClient.On("AcknowledgeMails", getTestAcknowledges()).Return(nil)
	mockedIO.On("Print", "E-Mails are Acknowledged: ")
	mockedStorage.On("DeleteAcknowledges", testAcknowledges).Return(nil)

	testee := CreateReciever(config, mockedIO, mockedApiClient, mockedStorage, mockedConfirm)
	testee.Run()

	mockedIO.AssertExpectations(t)
	mockedApiClient.AssertExpectations(t)
	mockedStorage.AssertExpectations(t)
	mockedConfirm.AssertExpectations(t)
}

func Test_AcknowledgeSpecify(t *testing.T) {
	testAcknowledges := getTestAcknowledges()
	testSpecifyAcks := getSpecifyAcknowledges()
	testSelectedAcks := getTestSelectedAcknowledge()
	testMails := getTestMails()
	config := configuration.Configuration{}
	mockedIO := new(inputOutput.MockedInputOutput)
	mockedApiClient := new(client.MockedClient)
	mockedConfirm := new(confirm.MockedConfirm)
	mockedStorage := new(acknowledgementStorage.MockedAcknowledgementStorage)

	mockedApiClient.On("ReceiveMails").Return(testMails, nil)
	mockedIO.On("Print", strconv.Itoa(len(testMails))+" Mails are coming from Server")
	mockedConfirm.On("GetAllAcknowledges", testMails).Return(testAcknowledges)
	mockedStorage.On("AppendAcknowledgements", testAcknowledges).Return(nil)
	mockedIO.On("Print", "Save Acknowledges...")
	mockedStorage.On("ReadAcknowledgements").Return(testAcknowledges, nil)
	mockedIO.On("Print", "Available Mails: "+strconv.Itoa(len(testAcknowledges)))
	mockedIO.On("Print", "send all Acknowledges or specify Acknowledges to Server. Or stop reciever (all/specify/stop):")
	mockedIO.On("ReadEntry").Return("specify").Once()
	mockedConfirm.On("ShowAllEmailAcks", testAcknowledges)
	mockedIO.On("Print", "Specify Acknowledge by Subject: ")
	mockedIO.On("ReadEntry").Return("testSubject1").Once()
	mockedIO.On("ReadEntry").Return("stop")
	mockedConfirm.On("GetSingleAcknowledges", testAcknowledges, "testSubject1").Return(testSpecifyAcks, testSelectedAcks)
	mockedApiClient.On("AcknowledgeMails", testSelectedAcks).Return(nil)
	mockedStorage.On("DeleteAcknowledges", testSelectedAcks).Return(nil)
	mockedIO.On("Print", "E-Mail is Acknowledged: ")

	testee := CreateReciever(config, mockedIO, mockedApiClient, mockedStorage, mockedConfirm)
	testee.Run()

	mockedIO.AssertExpectations(t)
	mockedApiClient.AssertExpectations(t)
	mockedStorage.AssertExpectations(t)
	mockedConfirm.AssertExpectations(t)
}

func Test_AcknowledgeStop(t *testing.T) {
	config := configuration.Configuration{}
	mockedIO := new(inputOutput.MockedInputOutput)
	mockedApiClient := new(client.MockedClient)
	mockedConfirm := new(confirm.MockedConfirm)
	mockedStorage := new(acknowledgementStorage.MockedAcknowledgementStorage)

	mockedApiClient.On("ReceiveMails").Return(getTestMails(), nil)
	mockedIO.On("Print", strconv.Itoa(len(getTestMails()))+" Mails are coming from Server")
	mockedConfirm.On("GetAllAcknowledges", getTestMails()).Return(getTestAcknowledges())
	mockedStorage.On("AppendAcknowledgements", getTestAcknowledges()).Return(nil)
	mockedIO.On("Print", "Save Acknowledges...")
	mockedStorage.On("ReadAcknowledgements").Return(getTestAcknowledges(), nil)
	mockedIO.On("Print", "Available Mails: "+strconv.Itoa(len(getTestAcknowledges())))
	mockedIO.On("Print", "send all Acknowledges or specify Acknowledges to Server. Or stop reciever (all/specify/stop):")
	mockedIO.On("ReadEntry").Return("stop")

	testee := CreateReciever(config, mockedIO, mockedApiClient, mockedStorage, mockedConfirm)
	testee.Run()

	mockedIO.AssertExpectations(t)
	mockedApiClient.AssertExpectations(t)
	mockedStorage.AssertExpectations(t)
	mockedConfirm.AssertExpectations(t)

}

func TestRecieveMailsError(t *testing.T) {

	testMails := getTestMails()

	config := configuration.Configuration{}
	mockedIO := new(inputOutput.MockedInputOutput)
	mockedApiClient := new(client.MockedClient)
	mockedConfirm := new(confirm.MockedConfirm)
	mockedStorage := new(acknowledgementStorage.MockedAcknowledgementStorage)

	mockedApiClient.On("ReceiveMails").Return(testMails, errors.New(""))
	mockedIO.On("Print", "Transmission is going wrong. Retry? (n,press any key)")
	mockedIO.On("ReadEntry").Return("n")

	testee := CreateReciever(config, mockedIO, mockedApiClient, mockedStorage, mockedConfirm)
	testee.Run()

	mockedIO.AssertExpectations(t)
	mockedApiClient.AssertExpectations(t)
	mockedStorage.AssertExpectations(t)
	mockedConfirm.AssertExpectations(t)
}
