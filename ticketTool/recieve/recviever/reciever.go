package recviever

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/acknowledgementStorage"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/confirm"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"fmt"
	"strconv"
)

type Reciever struct {
	io            inputOutput.InputOutput
	apiClient     client.ApiClient
	storage       acknowledgementStorage.AckStorage
	recieveConfig configuration.Configuration
}

func CreateReciever(config configuration.Configuration, io inputOutput.InputOutput, apiClient client.ApiClient, storage acknowledgementStorage.AckStorage) Reciever {
	reciever := Reciever{io: io, apiClient: apiClient, storage: storage, recieveConfig: config}
	return reciever
}

func (r *Reciever) Run() error {
	recieveMails, err := r.apiClient.ReceiveMails()
	if err != nil {
		r.io.Print("Transmission is going wrong. Retry? (n,press any key)")
		answer := r.io.ReadEntry()
		if answer == "n" {
			return err
		}
	} else {
		r.io.Print(strconv.Itoa(len(recieveMails)) + " Mails are coming from Server")
		acknowledges := confirm.GetAllAcknowledges(recieveMails)
		r.storage.AppendAcknowledgements(acknowledges)
		r.io.Print("Save Acknowledges...")
		allAcknowledges, err := r.storage.ReadAcknowledgements()
		if err != nil {
			r.io.Print("couldn't read storaged Acknowledges")
			return err
		} else if len(allAcknowledges) == 0 {
			r.io.Print("No Emails available")
			return nil
		}
		r.io.Print("Available Mails: " + strconv.Itoa(len(allAcknowledges)))
		if len(allAcknowledges) != 0 {
			r.allOrSpecifyConfirm(&allAcknowledges)
			return nil
		}
	}
	return nil
}

func (r *Reciever) allOrSpecifyConfirm(allAcknowledges *[]mail.Acknowledgment) {
	for true {
		fmt.Println("send all Acknowledges or specify Acknowledges to Server. Or stop reciever (all/specify/stop):")
		answer := r.io.ReadEntry()
		if answer == "all" {
			ackError := r.apiClient.AcknowledgeMails(*allAcknowledges)
			if ackError != nil {
				fmt.Println("acknowlege is not posted")
			} else {
				fmt.Println("E-Mails are Acknowledged: ")
				r.storage.DeleteAcknowledges(*allAcknowledges)
				break
			}
		} else if answer == "specify" {
			confirm.ShowAllEmailAcks(*allAcknowledges)
			fmt.Println("Specify Acknowledge by Subject: ")
			answer := r.io.ReadEntry()
			newAcknowledges, selectedAck := confirm.GetSingleAcknowledges(*allAcknowledges, answer)
			allAcknowledges = &newAcknowledges
			ackError := r.apiClient.AcknowledgeMails(selectedAck)
			if ackError != nil {
				fmt.Println("acknowlege is not posted")
			} else {
				r.storage.DeleteAcknowledges(selectedAck)
				fmt.Println("E-Mail is Acknowledged: ")
			}
			if len(*allAcknowledges) == 0 {
				break
			}
		} else if answer == "stop" {
			break
		}
	}
}
