package mail

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/validation/mail"
	"encoding/json"
	"github.com/pkg/errors"
	"html"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

/*
	Interface for the mail context.
*/
type MailContext interface {
	GetUnsentMails() ([]Mail, error)
	AcknowledgeMails(acknowledgments []Acknowledgment) error
	CreateNewOutgoingMail(receiver string, subject string, content string) error
}

/*
	A mail manager.
*/
type MailManager struct {
	unAcknowledgedMails     []string
	unAcknowledgeMailsMutex sync.RWMutex
	unSentMails             []Mail
	unSentMailsMutex        sync.RWMutex
	mailFolderPath          string
	mailFolderAccessMutex   sync.RWMutex
	logger                  logging.Logger
	outgoingMailAddress     string
}

/*
	Acknowledging a mail deletes it.
*/
func (t *MailManager) AcknowledgeMails(acknowledgments []Acknowledgment) error {
	t.unAcknowledgeMailsMutex.Lock()
	defer t.unAcknowledgeMailsMutex.Unlock()

	for _, acknowledgment := range acknowledgments {
		t.unAcknowledgedMails = remove(t.unAcknowledgedMails, acknowledgment.Id)
		mailPath := path.Join(t.mailFolderPath, acknowledgment.Id+".json")
		exist, err := helpers.FilePathExists(mailPath)
		if err != nil {
			return err
		}
		if exist {
			err := os.Remove(mailPath)
			if err != nil {
				return err
			}
		}
	}
	err := t.persistUnAcknowledgedMailState()
	if err != nil {
		return err
	}
	return nil
}

/*
	Remove a string value from a string array.
*/
func remove(array []string, valueRoRemove string) []string {
	for i, v := range array {
		if v == valueRoRemove {
			array = append(array[:i], array[i+1:]...)
			break
		}
	}
	return array
}

/*
	Get the unsent emails.
*/
func (t *MailManager) GetUnsentMails() ([]Mail, error) {
	t.unSentMailsMutex.Lock()
	defer t.unSentMailsMutex.Unlock()
	t.unAcknowledgeMailsMutex.Lock()
	defer t.unAcknowledgeMailsMutex.Unlock()

	mailsToSent := t.unSentMails
	for _, mailToSent := range t.unSentMails {
		t.unAcknowledgedMails = append(t.unAcknowledgedMails, mailToSent.Id)
	}
	err := t.persistUnAcknowledgedMailState()
	if err != nil {
		return []Mail{}, err
	}
	t.unSentMails = []Mail{}
	return mailsToSent, nil

}

/*
	Initialize the MailManager with the given folder path. The folder is used to load and store the mail data.
*/
func (t *MailManager) Initialize(folderPath string, outgoingMailAddress string, logger logging.Logger) error {
	t.logger = logger
	if folderPath == "" {
		return errors.New("path to mail data storage can not be a empty string.")
	}
	validator := mail.NewValidator()
	if !validator.Validate(outgoingMailAddress) {
		return errors.New("Outgoing mail address must be valid.")
	}
	t.outgoingMailAddress = outgoingMailAddress

	t.unAcknowledgedMails = []string{}
	t.unSentMails = []Mail{}
	t.mailFolderPath = folderPath
	folderExists, err := helpers.FilePathExists(folderPath)
	if err != nil {
		return errors.Wrap(err, "could not check if folder for mails exists.")
	}
	if !folderExists {
		t.logger.LogInfo("MailManager", "Mail data folder does not exist, going to create it")
		err = helpers.CreateFolderPath(folderPath)
		if err != nil {
			return errors.Wrap(err, "could not create folder for mails.")
		}
	} else {
		return t.readExistingMailData()
	}
	return nil
}

/*
	Create a new outgoing mail
*/
func (t *MailManager) CreateNewOutgoingMail(receiver string, subject string, content string) error {
	validator := mail.NewValidator()
	if !validator.Validate(receiver) {
		return errors.New("Receiver address is not valid")
	}
	adjustedContent := html.EscapeString(content)
	adjustedSubject := html.EscapeString(subject)
	uuid, err := helpers.GenerateUUID()
	if err != nil {
		return err
	}
	mailToSent := Mail{Id: uuid, Sender: t.outgoingMailAddress, Receiver: receiver,
		Subject: adjustedSubject, Content: adjustedContent, SentTime: time.Now().Unix()}

	t.mailFolderAccessMutex.Lock()
	defer t.mailFolderAccessMutex.Unlock()
	err = t.persistMailToDisk(mailToSent)
	if err != nil {
		return err
	}
	t.unSentMailsMutex.Lock()
	defer t.unSentMailsMutex.Unlock()
	t.unSentMails = append(t.unSentMails, mailToSent)

	t.logger.LogInfo("MailManager", "Mail created with id: "+mailToSent.Id)
	return nil
}

/*
	Read existing mails from the file system.
*/
func (t *MailManager) readExistingMailData() error {
	t.mailFolderAccessMutex.Lock()
	defer t.mailFolderAccessMutex.Unlock()

	err := t.readUnAcknowledgeMailsFile(t.mailFolderPath)
	if err != nil {
		return errors.Wrap(err, "could not read existing mail data.")
	}

	t.unSentMailsMutex.Lock()
	t.unSentMailsMutex.Unlock()

	// Read existing mail files:
	files, err := ioutil.ReadDir(t.mailFolderPath)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		// Ignore folders
		if !fileInfo.IsDir() {
			match, _ := regexp.MatchString(".json", fileInfo.Name())
			if match {
				if fileInfo.Name() == unAcknowledgeMailFileName {
					continue
				}
				extension := filepath.Ext(fileInfo.Name())
				name := fileInfo.Name()[0 : len(fileInfo.Name())-len(extension)]
				ignoreFile := false
				for _, alreadySentId := range t.unAcknowledgedMails {
					if name == alreadySentId {
						ignoreFile = true
						break
					}
				}
				if !ignoreFile {
					mail, err := t.readMailFromFile(path.Join(t.mailFolderPath, fileInfo.Name()))
					if err != nil {
						t.logger.LogInfo("MailManager", "Could not decode data from mail file. Going to ignore it. See error below:")
						t.logger.LogError("MailManager", err)
						continue
					}
					t.unSentMails = append(t.unSentMails, *mail)
				} else {
					t.logger.LogInfo("MailManager", "Not loading mail file "+fileInfo.Name()+" since it is waiting for acknowledgement")
				}
			}
		}
	}
	return nil
}

/*
	Read and decode a mail from a file.
*/
func (t *MailManager) readMailFromFile(pathToFile string) (*Mail, error) {
	t.logger.LogInfo("MailManager", "Loading mail file "+pathToFile)

	fileValue, err := helpers.ReadAllDataFromFile(pathToFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not read data from file")
	}

	parsedData := new(Mail)
	err = json.Unmarshal(fileValue, &parsedData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode data from file")
	}
	return parsedData, nil
}

/*
	Read the list of unacknowledged mails.
*/
func (t *MailManager) readUnAcknowledgeMailsFile(folderPath string) error {
	t.unAcknowledgeMailsMutex.Lock()
	defer t.unAcknowledgeMailsMutex.Unlock()

	filePath := path.Join(folderPath, unAcknowledgeMailFileName)
	err := helpers.CreateFileIfNotExists(filePath)
	if err != nil {
		return err
	}
	// Read data from file
	fileValue, err := helpers.ReadAllDataFromFile(filePath)
	if err != nil {
		return err
	}

	parsedData := new([]string)
	json.Unmarshal(fileValue, &parsedData)
	if err != nil {
		t.logger.LogInfo("MailManager", "Could not decode data from file, possibly empty. Create new file")
		jsonData, err := json.Marshal([]string{})
		if err != nil {
			return errors.Wrap(err, "could not create empty data file for unAcknowledged mails.")
		}
		helpers.WriteDataToFile(filePath, jsonData)
		if err != nil {
			return errors.Wrap(err, "could not write data file for unAcknowledged mails.")
		}
	}
	t.unAcknowledgedMails = *parsedData
	return nil
}

/*
	Persist a mail to disk.
*/
func (t *MailManager) persistMailToDisk(mailToSent Mail) error {
	jsonData, err := json.MarshalIndent(mailToSent, "", "    ")
	if err != nil {
		return errors.Wrap(err, "could not save mail to file")
	}
	return helpers.WriteDataToFile(path.Join(t.mailFolderPath, mailToSent.Id+".json"), jsonData)
}

/*
	Persist the unacknowledged state to disk.
*/
func (t *MailManager) persistUnAcknowledgedMailState() error {
	jsonData, err := json.Marshal(t.unAcknowledgedMails)
	if err != nil {
		return errors.Wrap(err, "could not encode data file for unAcknowledged mails.")
	}
	helpers.WriteDataToFile(path.Join(t.mailFolderPath, unAcknowledgeMailFileName), jsonData)
	if err != nil {
		return errors.Wrap(err, "could not write data file for unAcknowledged mails.")
	}
	return nil
}

const unAcknowledgeMailFileName = "unAcknowledgedMails.json"
