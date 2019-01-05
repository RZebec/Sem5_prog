package recviever

import (
	"de/vorlesung/projekt/IIIDDD/ticketTool/client"
	"de/vorlesung/projekt/IIIDDD/ticketTool/configuration"
	"de/vorlesung/projekt/IIIDDD/ticketTool/inputOutput"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/acknowledgementStorage"
	"de/vorlesung/projekt/IIIDDD/ticketTool/recieve/confirm"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"fmt"
	"strconv"
)

type Reciever struct {
	io            inputOutput.InputOutput
	apiClient     client.Client
	storage       acknowledgementStorage.AckStorage
	recieveConfig configuration.Configuration
	confirmator   confirm.Confirmation
}

/*
initialize Reciever
*/
func CreateReciever(config configuration.Configuration, io inputOutput.InputOutput, apiClient client.Client,
	storage acknowledgementStorage.AckStorage, confirmation confirm.Confirmation) Reciever {
	reciever := Reciever{io: io, apiClient: apiClient, storage: storage, recieveConfig: config, confirmator: confirmation}
	return reciever
}

func (r *Reciever) Run() error {
	recieveMails, err := r.apiClient.ReceiveMails() //recieve Mails from Server
	if err != nil {
		r.io.Print("Transmission is going wrong. Retry? (n,press any key)")
		answer := r.io.ReadEntry() //request if you want to retry the transmission
		if answer == "n" {
			return err
		}
	} else {
		r.io.Print(strconv.Itoa(len(recieveMails)) + " Mails are coming from Server")
		for _, receivedMail := range recieveMails {
			r.io.Print("Receiver: " + receivedMail.Receiver + " Subject: " + receivedMail.Subject)
		}
		acknowledges := r.confirmator.GetAllAcknowledges(recieveMails) //create Acknowledges from Mails
		err := r.storage.AppendAcknowledgements(acknowledges)          //store these Acknowledges
		if err != nil {
			r.io.Print("mails cant't saved: " + err.Error())
			return err
		}
		r.io.Print("Save Acknowledges...")
		allAcknowledges, err := r.storage.ReadAcknowledgements() //read the recieved Acknowledges and previous Acknowledges
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

/*
you can select between one specify confirm or if you want all confirm
*/
func (r *Reciever) allOrSpecifyConfirm(allAcknowledges *[]mailData.Acknowledgment) {
	for true {
		r.io.Print("send all Acknowledges or specify Acknowledges to Server. Or stop reciever (all/specify/stop):")
		answer := r.io.ReadEntry()
		if answer == "all" {
			ackError := r.apiClient.AcknowledgeMails(*allAcknowledges) //send a list from all Acknowledges back to Server
			if ackError != nil {
				r.io.Print("acknowlege is not posted")
			} else {
				r.io.Print("E-Mails are Acknowledged: ")
				err := r.storage.DeleteAcknowledges(*allAcknowledges) //if the send process was sucessfull, delete all Acknowledges from your Storage
				if err != nil {
					r.io.Print("Acknowledges couldn't deleted: " + err.Error())
					break
				}
				break
			}
		} else if answer == "specify" { // specify confirmation
			r.confirmator.ShowAllEmailAcks(*allAcknowledges)
			r.io.Print("Specify Acknowledge by Subject: ")
			answer := r.io.ReadEntry()
			/*
				get back a List with one Acknowledge and delete the selected Acknowledge from the List of all
			*/
			newAcknowledges, selectedAck := r.confirmator.GetSingleAcknowledges(*allAcknowledges, answer)
			allAcknowledges = &newAcknowledges                    //all Acknowledges - selected Acknowledge = allAcknowledges
			ackError := r.apiClient.AcknowledgeMails(selectedAck) //send back the List with one Acknowledge
			if ackError != nil {
				fmt.Println("acknowlege is not posted")
			} else {
				err := r.storage.DeleteAcknowledges(selectedAck) //if send process was sucessfull, delete the selected Acknowledge
				if err != nil {
					r.io.Print("Selected Acknowledge couldn't deleted: " + err.Error())
				}
				r.io.Print("E-Mail is Acknowledged: ")
			}
			if len(*allAcknowledges) == 0 {
				break
			}
		} else if answer == "stop" {
			break
		}
	}
}
