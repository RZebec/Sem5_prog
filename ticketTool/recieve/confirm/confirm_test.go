package confirm

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestGetAllAcknowledges(t *testing.T) {
	mails := []mail.Mail{
		{Id: "123", Sender: "abc@gmx.de", Receiver: "defg@web.de", Subject: "testSubject0", Content: "Test", SentTime: time.Now().Unix()},
		{Id: "456", Sender: "abc@gmx.de", Receiver: "defg@web.de", Subject: "testSubject1", Content: "Test", SentTime: time.Now().Unix()},
		{Id: "789", Sender: "abc@gmx.de", Receiver: "defg@web.de", Subject: "testSubject2", Content: "Test", SentTime: time.Now().Unix()}}

	acknowledge := GetAllAcknowledges(mails)

	assert.True(t, acknowledge[0].Subject == mails[0].Subject, "Subject is the same")
	assert.True(t, acknowledge[0].Id == mails[0].Id, "Id is the same")

	assert.True(t, acknowledge[1].Subject == mails[1].Subject, "Subject is the same")
	assert.True(t, acknowledge[1].Id == mails[1].Id, "Id is the same")

	assert.True(t, acknowledge[2].Subject == mails[2].Subject, "Subject is the same")
	assert.True(t, acknowledge[2].Id == mails[2].Id, "Id is the same")
}

func TestGetSingleAcknowledges(t *testing.T) {
	mails := []mail.Mail{
		{Id: "123", Sender: "abc@gmx.de", Receiver: "defg@web.de", Subject: "testSubject0", Content: "Test", SentTime: time.Now().Unix()},
		{Id: "456", Sender: "abc@gmx.de", Receiver: "defg@web.de", Subject: "testSubject1", Content: "Test", SentTime: time.Now().Unix()},
		{Id: "789", Sender: "abc@gmx.de", Receiver: "defg@web.de", Subject: "testSubject2", Content: "Test", SentTime: time.Now().Unix()}}

	acknowledges := GetAllAcknowledges(mails)
	/*
		selected acknowledge should be not anymore in the List of acknowledges
	*/
	acknowledgeList, selectedAck := GetSingleAcknowledges(acknowledges, "testSubject0")
	for i := 0; i < len(acknowledgeList); i++ {
		assert.False(t, acknowledgeList[i].Id == selectedAck[0].Id, "selected Acknowledge is in acknowledge List")
		assert.True(t, acknowledgeList[i].Id != selectedAck[0].Id, "selected Acknowledge is not in acknowledge List")
	}

	/*
		if subject not exist, you selected Acknowledge List have length 0
		and the original acknowledge List is not changed
	*/
	equalAcknowledgeList := GetAllAcknowledges(mails)
	acknowledgeList, selectedAck = GetSingleAcknowledges(acknowledges, "testSubject3")
	assert.True(t, len(selectedAck) == 0, strconv.Itoa(len(selectedAck))+" :should be 0")

	assert.Equal(t, len(equalAcknowledgeList), len(acknowledgeList), "length of list shouldnt be change")
}

func TestDeleteFromArray(t *testing.T) {
	mails := []mail.Mail{
		{Id: "123", Sender: "abc@gmx.de", Receiver: "defg@web.de", Subject: "testSubject0", Content: "Test", SentTime: time.Now().Unix()},
		{Id: "456", Sender: "abc@gmx.de", Receiver: "defg@web.de", Subject: "testSubject1", Content: "Test", SentTime: time.Now().Unix()},
		{Id: "789", Sender: "abc@gmx.de", Receiver: "defg@web.de", Subject: "testSubject2", Content: "Test", SentTime: time.Now().Unix()}}

	acknowledges := GetAllAcknowledges(mails)

	newAcknowledges := deleteFromArray(acknowledges, acknowledges[0])
	assert.Equal(t, newAcknowledges[0].Id, "456", "Indexing is wrong")
	assert.Equal(t, newAcknowledges[1].Id, "789", "Indexing is wrong")
	assert.True(t, len(acknowledges)-1 == len(newAcknowledges), "element is not deleted")
}
