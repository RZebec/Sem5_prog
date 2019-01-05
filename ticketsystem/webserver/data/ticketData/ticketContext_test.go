package ticketData

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/core/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

/*
	Example for the initialization of the ticketData manager.
*/
func ExampleTicketManager_Initialize() {
	// Just to clean up after the example:
	defer os.RemoveAll("/temp/pathToTicketFolder")
	ticketContext := TicketManager{}
	err := ticketContext.Initialize("/temp/pathToTicketFolder")
	fmt.Println(err)
	// Output:
	// <nil>
}

/*
	Example to get a ticketData by its id.
*/
func ExampleTicketManager_GetTicketById() {
	ticketContext := TicketManager{}
	initializeWithTempTicketForExample(&ticketContext)

	exists, ticket := ticketContext.GetTicketById(1)
	fmt.Println(exists)
	fmt.Println(ticket.Info().Id)
	// Output:
	// true
	// 1
}

/*
	Example to get all ticketData infos.
*/
func ExampleTicketManager_GetAllTicketInfo() {
	ticketContext := TicketManager{}
	initializeWithTempTicketForExample(&ticketContext)

	ticket := ticketContext.GetAllTicketInfo()
	fmt.Println(ticket[0].Id)
	fmt.Println(ticket[0].Title)
	// Output:
	// 1
	// Example_TicketTitle
}

/*
	Example to create a new ticketData for an internal userData.
*/
func ExampleTicketManager_CreateNewTicketForInternalUser() {
	// Preparation for example:
	folderPath, rootPath, _ := prepareTempDirectory()
	defer os.RemoveAll(rootPath)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	createdUser := userData.User{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: createdUser.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicketForInternalUser("Example_TestTitle", createdUser, initialMessage)

	fmt.Println(err)
	fmt.Println(newTicket.Info().Title)
	// Output:
	// <nil>
	// Example_TestTitle
}

/*
	Example to create a new ticketData.
*/
func ExampleTicketManager_CreateNewTicket() {
	// Preparation for example:
	folderPath, rootPath, _ := prepareTempDirectory()
	defer os.RemoveAll(rootPath)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	creator := Creator{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicket("Example_TestTitle", creator, initialMessage)

	fmt.Println(err)
	fmt.Println(newTicket.Info().Title)
	// Output:
	// <nil>
	// Example_TestTitle
}

/*
	Example to append a new message to a ticketData.
*/
func ExampleTicketManager_AppendMessageToTicket() {
	ticketContext := TicketManager{}
	folderPath := initializeWithTicketForExample(&ticketContext)
	defer os.RemoveAll(folderPath)

	exists, ticket := ticketContext.GetTicketById(1)
	fmt.Println(exists)
	fmt.Println(len(ticket.Messages()))

	newMessage := MessageEntry{CreatorMail: "alex@wagner.de", Content: "This is the message", OnlyInternal: false}

	updatedTicket, _ := ticketContext.AppendMessageToTicket(ticket.Info().Id, newMessage)

	fmt.Println(len(updatedTicket.Messages()))
	fmt.Println(updatedTicket.Messages()[1].Content)

	// Output:
	// true
	// 1
	// 2
	// This is the message
}

/*
	Setting the editor when no editor is set should work.
*/
func TestTicketManager_SetEditor_NoPreviousEditor_EditorIsSet(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	// Ensure, that the ticketData has no editor set:
	_, testTicket := testee.GetTicketById(2)
	testEditor := userData.User{Mail: "test@test", UserId: 2, FirstName: "first", LastName: "last", Role: userData.RegisteredUser, State: userData.Active}
	assert.False(t, testTicket.info.HasEditor, "No editor should be set")

	// Set the editor:
	testee.SetEditor(testEditor, testTicket.info.Id)

	// Assert that the ticketData has been updated:
	_, updatedTicket := testee.GetTicketById(testTicket.info.Id)
	assert.True(t, updatedTicket.info.HasEditor, "Editor should be set")
	assert.Equal(t, 2, updatedTicket.info.Editor.UserId, "Editor should be set")
	assert.Equal(t, testEditor, updatedTicket.info.Editor, "Editor should be set")

	// Assert that the stored file has been updated:
	storedTicket, err := readTicketFromFile(updatedTicket.filePath)
	assert.Equal(t, testEditor, storedTicket.info.Editor)
}

/*
	Setting the editor when no editor is set should work.
*/
func TestTicketManager_SetEditor_PreviousEditorSet_EditorIsUpdated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	// Ensure, that the ticketData has a editor set:
	_, testTicket := testee.GetTicketById(3)
	testEditor := userData.User{Mail: "test@test", UserId: 2, FirstName: "first", LastName: "last", Role: userData.RegisteredUser, State: userData.Active}
	assert.True(t, testTicket.info.HasEditor, "Editor should be set")

	// Update the editor:
	_, err = testee.SetEditor(testEditor, testTicket.info.Id)
	assert.Nil(t, err)

	// Assert that the ticketData has been updated:
	_, updatedTicket := testee.GetTicketById(testTicket.info.Id)
	assert.True(t, updatedTicket.info.HasEditor, "Editor should be set")
	assert.Equal(t, 2, updatedTicket.info.Editor.UserId, "Editor should be set")
	assert.Equal(t, testEditor, updatedTicket.info.Editor, "Editor should be set")

	// Assert that the stored file has been updated:
	storedTicket, err := readTicketFromFile(updatedTicket.filePath)
	assert.Equal(t, testEditor, storedTicket.info.Editor)
}

/*
	Removing a editor from a ticketData should set the invalid default editor.
*/
func TestTicketManager_RemoveEditor(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	// Ensure, that the ticketData has a editor set:
	_, testTicket := testee.GetTicketById(3)
	assert.True(t, testTicket.info.HasEditor, "Editor should be set")

	// Remove the editor:
	err = testee.RemoveEditor(testTicket.info.Id)
	assert.Nil(t, err)

	// Assert that the ticketData has been updated:
	_, updatedTicket := testee.GetTicketById(testTicket.info.Id)
	assert.False(t, updatedTicket.info.HasEditor, "Editor should not be set")
	assert.Equal(t, userData.GetInvalidDefaultUser(), updatedTicket.info.Editor, "Editor should be set to invalid id 0")

	// Assert that the stored file has been updated:
	storedTicket, err := readTicketFromFile(updatedTicket.filePath)
	assert.Equal(t, userData.GetInvalidDefaultUser(), storedTicket.info.Editor)
	assert.False(t, storedTicket.info.HasEditor)
}

/*
	Removing a editor from a non existing ticketData should return a error.
*/
func TestTicketManager_RemoveEditor_TicketDoesNotExist(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	// Remove the editor:
	err = testee.RemoveEditor(9999)
	assert.Equal(t, "ticketData does not exist", err.Error())
}

/*
	A ticketData can be merged with another ticketData. The messages from the newer tickets will be attached to the older ticketData.
	The newer ticketData will be deleted.
*/
func TestTicketManager_MergeTickets_TicketsAreMerged(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	_, firstTicket := testee.GetTicketById(1)
	_, secondTicket := testee.GetTicketById(2)

	// Tickets need to have the same editor:
	testEditor := userData.User{Mail: "test@test", UserId: 2, FirstName: "first", LastName: "last", Role: userData.RegisteredUser, State: userData.Active}
	_, err = testee.SetEditor(testEditor, 1)
	assert.Nil(t, err)
	_, err = testee.SetEditor(testEditor, 2)
	assert.Nil(t, err)

	// Merge the tickets:
	success, err := testee.MergeTickets(1, 2)
	assert.True(t, success, "Tickets should be merged")
	assert.Nil(t, err)

	// Ticket 1 is older, so Ticket 2 should be deleted:
	exists, _ := testee.GetTicketById(2)
	assert.False(t, exists, "Ticket 2 should be deleted")

	exists, mergedTicket := testee.GetTicketById(1)
	assert.True(t, exists, "Ticket 1 should still exist")

	// Assert that the messages have been merged:
	assert.Equal(t, len(firstTicket.messages)+len(secondTicket.messages), len(mergedTicket.messages),
		"Merged ticketData should contain all messages")

	// All messages from the first ticketData should be merged:
	for _, message := range firstTicket.messages {
		found := false
		for _, mergedMessage := range mergedTicket.messages {
			if message.Content == mergedMessage.Content {
				found = true
				break
			}
		}
		assert.True(t, found, "Messages from the first ticketData should be merged")
	}
	// All messages from the second ticketData should be merged:
	for _, message := range secondTicket.messages {
		found := false
		for _, mergedMessage := range mergedTicket.messages {
			if message.Content == mergedMessage.Content {
				found = true
				break
			}
		}
		assert.True(t, found, "Messages from the second ticketData should be merged")
	}
	secondTicketFileExists, err := helpers.FilePathExists(secondTicket.filePath)
	assert.False(t, secondTicketFileExists, "The file for the second ticketData should be deleted")
	assert.Nil(t, err)

	cached, _ := testee.GetTicketById(secondTicket.info.Id)
	assert.False(t, cached, "The second ticketData should be removed from the cache")
}

/*
	Merging a non existing ticketData should not be possible. A error should be returned.
*/
func TestTicketManager_MergeTickets_TicketDoesNotExist(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	success, err := testee.MergeTickets(9999, 1)
	assert.False(t, success, "Merging a non existing ticketData should not be possible")
	assert.Equal(t, "ticketData not found", err.Error())

	success, err = testee.MergeTickets(1, 9999)
	assert.False(t, success, "Merging a non existing ticketData should not be possible")
	assert.Equal(t, "ticketData not found", err.Error())
}

/*
	Merging a ticketData without a editor should not be possible. A error should be returned.
*/
func TestTicketManager_MergeTickets_NoEditorSet(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	success, err := testee.MergeTickets(1, 2)
	assert.False(t, success, "Merging a ticketData without editor should not be possible")
	assert.Equal(t, "can not merge ticketData if there is no editor", err.Error())

	success, err = testee.MergeTickets(2, 1)
	assert.False(t, success, "Merging a ticketData without editor should not be possible")
	assert.Equal(t, "can not merge ticketData if there is no editor", err.Error())
}

/*
	Merging a tickets with different editors should not be possible. A error should be returned.
*/
func TestTicketManager_MergeTickets_DifferentEditorSet(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	success, err := testee.MergeTickets(3, 4)
	assert.False(t, success, "Merging a tickets with different editors should not be possible")
	assert.Equal(t, "only tickets of the same editor can be merged", err.Error())

	success, err = testee.MergeTickets(4, 3)
	assert.False(t, success, "Merging a tickets with different editors should not be possible")
	assert.Equal(t, "only tickets of the same editor can be merged", err.Error())
}

/*
	Merging a tickets with the same id should not be possible. A error should be returned.
*/
func TestTicketManager_MergeTickets_MergeTicketWithItself(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	success, err := testee.MergeTickets(1, 1)
	assert.False(t, success, "Merging a ticketData with itself should not be possible")
	assert.Equal(t, "can not merge a ticketData with itself", err.Error())
}

/*
	Getting the ticketData infos when no tickets exists, should return a empty array.
*/
func TestTicketManager_GetAllTicketInfo_NoTickets_EmptyArrayReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetAllTicketInfo()
	assert.Equal(t, 0, len(tickets))
}

/*
	Getting the open ticketData infos when ticketData exist, should return the infos for the open tickets.
*/
func TestTicketManager_GetAllOpenTickets_TicketsExist_TicketInfoReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetAllOpenTickets()
	assert.Equal(t, 2, len(tickets))
}

/*
	Getting the ticketData infos for open tickets when no tickets exists, should return a empty array.
*/
func TestTicketManager_GetAllOpenTickets_NoTickets_EmptyArrayReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetAllOpenTickets()
	assert.Equal(t, 0, len(tickets))
}

/*
	Where there are tickets which where created by the given mail, they should be returned.
 */
func TestTicketManager_GetTicketsForCreatorMail_TicketExists_TicketInfoReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetTicketsForCreatorMail("test@test.de")
	assert.Equal(t, 3, len(tickets))
}

/*
	Where there are no tickets which where created by the given mail, an empty array should be returned.
 */
func TestTicketManager_GetTicketsForCreatorMail_NoTicketExists_EmptyTicketInfoReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetTicketsForCreatorMail("abc@temp.de")
	assert.Equal(t, 0, len(tickets))
}

/*
	If there are tickets with the given id as editor, they should be returned.
 */
func TestTicketManager_GetTicketsForEditorId_TicketExists_TicketsReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetTicketsForEditorId(2)
	assert.Equal(t, 1, len(tickets))
}

/*
	If there are no tickets with the given id as editor, a empty array should be returned.
 */
func TestTicketManager_GetTicketsForEditorId_NoTicketExists_NoTicketsReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetTicketsForEditorId(9)
	assert.Equal(t, 0, len(tickets))
}

/*
	Getting the ticketData infos when ticketData exist, should return the infos for the existing tickets.
*/
func TestTicketManager_GetAllTicketInfo_TicketsExist_TicketInfoReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	tickets := testee.GetAllTicketInfo()
	assert.Equal(t, 4, len(tickets))

	// Assert, that the correct data is returned:
	for _, v := range tickets {
		// Compare ticketData 4:
		if v.Id == 4 {
			assert.Equal(t, "TestTitle4", v.Title)
			assert.Equal(t, "peter@test.de", v.Editor.Mail)
			assert.Equal(t, 2, v.Editor.UserId)
			assert.Equal(t, "Peter", v.Editor.FirstName)
			assert.Equal(t, "Test", v.Editor.LastName)
			assert.Equal(t, true, v.HasEditor)
			assert.Equal(t, "peter@test.de", v.Creator.Mail)
			assert.Equal(t, "Peter", v.Creator.FirstName)
			assert.Equal(t, "Test", v.Creator.LastName)
			expectedCreationTime, _ := time.Parse(time.RFC3339Nano, "2018-10-27T18:23:04.1141357+02:00")
			assert.Equal(t, expectedCreationTime, v.CreationTime)
			expectedLastModificationTime, _ := time.Parse(time.RFC3339Nano, "2018-10-27T18:23:04.1141357+02:00")
			assert.Equal(t, expectedLastModificationTime, v.LastModificationTime)
		}
	}
}

/*
	Creating tickets for internal users in a concurrent way, should create all tickets.
*/
func TestTicketManager_CreateNewTicketsForInternalUser_ConcurrentAccess_AllCreated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)
	numberOfCreatedTickets := 100

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(numberOfCreatedTickets)
	for i := 0; i < numberOfCreatedTickets; i++ {

		go func(number int) {
			defer waitGroup.Done()
			id := strconv.Itoa(number)
			createdUser := userData.User{Mail: id + "@web.de", UserId: number, FirstName: "firstName" + id, LastName: "lastName" + id}
			initialMessage := MessageEntry{Id: 45, CreatorMail: createdUser.Mail, Content: "This is a test" + id, OnlyInternal: false}
			testee.CreateNewTicketForInternalUser("newTestTitle"+id, createdUser, initialMessage)

		}(i)
	}
	waitGroup.Wait()
	tickets := testee.GetAllTicketInfo()

	assert.Equal(t, numberOfCreatedTickets, len(tickets), "all ticketData info should be cached")
	// Check if ticketData data is correct
	for i := 0; i < numberOfCreatedTickets; i++ {
		id := strconv.Itoa(i)
		expectedCreator := Creator{Mail: id + "@web.de", FirstName: "firstName" + id, LastName: "lastName" + id}
		expectedTitle := "newTestTitle" + id
		found := false

		for _, ticket := range tickets {
			if ticket.Title == expectedTitle {
				found = true
				assert.Equal(t, expectedCreator, ticket.Creator, "the creator should be set")
			}
		}

		assert.True(t, found, "ticketData has not been found")
	}
}

/*
	Creating tickets for external users in a concurrent way, should create all tickets.
*/
func TestTicketManager_CreateNewTickets_ConcurrentAccess_AllCreated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)
	numberOfCreatedTickets := 100

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(numberOfCreatedTickets)
	for i := 0; i < numberOfCreatedTickets; i++ {

		go func(number int) {
			defer waitGroup.Done()
			id := strconv.Itoa(number)
			creator := Creator{Mail: id + "@web.de", FirstName: "Alex" + id, LastName: "Wagner" + id}
			initialMessage := MessageEntry{Id: 45, CreatorMail: creator.Mail, Content: "This is a test" + id, OnlyInternal: false}
			testee.CreateNewTicket("newTestTitle"+id, creator, initialMessage)

		}(i)
	}
	waitGroup.Wait()
	tickets := testee.GetAllTicketInfo()

	assert.Equal(t, numberOfCreatedTickets, len(tickets), "all ticketData info should be cached")
	// Check if ticketData data is correct
	for i := 0; i < numberOfCreatedTickets; i++ {
		id := strconv.Itoa(i)
		expectedCreator := Creator{Mail: id + "@web.de", FirstName: "Alex" + id, LastName: "Wagner" + id}
		expectedTitle := "newTestTitle" + id
		found := false

		for _, ticket := range tickets {
			if ticket.Title == expectedTitle {
				found = true
				assert.Equal(t, expectedCreator, ticket.Creator, "the creator should be set")
			}
		}

		assert.True(t, found, "ticketData has not been found")
	}
}

/*
	Creating a ticketData should really create it.
*/
func TestTicketManager_CreateNewTicket_TicketCreated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	creator := Creator{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 45, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicket("newTestTitle", creator, initialMessage)

	exists, createdTicket := testee.GetTicketById(newTicket.info.Id)
	assert.True(t, exists, "the ticketData should be created")

	// Validate that the file is created:
	exists, _ = helpers.FilePathExists(createdTicket.filePath)
	assert.True(t, exists, "the ticketData file should be created")

	storedTicket, err := readTicketFromFile(createdTicket.filePath)

	assertTicket(t, newTicket, storedTicket)
}

/*
	Creating a ticketData for a internal userData should really create the ticketData.
*/
func TestTicketManager_CreateNewTicketForInternalUser_TicketCreated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	createdUser := userData.User{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: createdUser.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicketForInternalUser("newTestTitle", createdUser, initialMessage)

	exists, createdTicket := testee.GetTicketById(newTicket.info.Id)
	assert.True(t, exists, "the ticketData should be created")

	// Validate that the file is created:
	exists, _ = helpers.FilePathExists(createdTicket.filePath)
	assert.True(t, exists, "the ticketData file should be created")

	storedTicket, err := readTicketFromFile(createdTicket.filePath)

	assertTicket(t, newTicket, storedTicket)
}

/*
	Initializing the ticketData manager with existing tickets, should load the ticketData.
*/
func TestTicketManager_Initialize_TicketsExist_TicketsAreLoaded(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	assert.Equal(t, 4, len(testee.cachedTickets))
	assert.Equal(t, 4, len(testee.cachedTicketIds))
}

/*
	Initializing the ticketData manager when no ticketData exist, should initialize without problems.
*/
func TestTicketManager_Initialize_NoTicketsExist_Initialized(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	assert.Equal(t, 0, len(testee.cachedTickets))
	assert.Equal(t, 0, len(testee.cachedTicketIds))
}

/*
	Initializing the ticketData manager with a invalid path, should return an error.
*/
func TestTicketManager_Initialize_InvalidFolderPath(t *testing.T) {
	testee := TicketManager{}
	err := testee.Initialize("")
	assert.Error(t, err, "folderPath is invalid")
}

/*
	Appending a message to a ticketData should append the message to the ticketData.
*/
func TestTicketManager_AppendMessageToTicket_MessageAppended(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	creator := Creator{Mail: "test1234@web.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 45, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	newTicket, err := testee.CreateNewTicket("newTestTitle", creator, initialMessage)

	// Id should be set by the context:
	message := MessageEntry{Id: 9999, CreatorMail: "max@muster.de", CreationTime: time.Now(),
		Content: "This is a appended message", OnlyInternal: false}
	ticket, err := testee.AppendMessageToTicket(newTicket.Info().Id, message)
	assert.Nil(t, err)

	found, storedTicket := testee.GetTicketById(newTicket.info.Id)
	assert.True(t, found)
	assert.Equal(t, ticket, storedTicket)

	for i, message := range storedTicket.messages {
		assert.Equal(t, i, message.Id, "the id should be in order")
	}
}

/*
	Appending a message to a non existing ticketData should return a error.
*/
func TestTicketManager_AppendMessageToTicket_TicketDoesNotExist_ErrorReturned(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	message := MessageEntry{Id: 9999, CreatorMail: "max@muster.de", CreationTime: time.Now(),
		Content: "This is a appended message", OnlyInternal: false}
	_, err = testee.AppendMessageToTicket(999, message)
	assert.Equal(t, "ticketData does not exist", err.Error())

}

/*
	Setting the state of the ticketData should change the state in memory and in the persisted ticketData.
*/
func TestTicketManager_SetTicketState_StateUpdated(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	// Ensure, that the ticketData has state open set:
	_, testTicket := testee.GetTicketById(2)
	assert.Equal(t, Open, testTicket.info.State, "State should be set to open")

	// Update the state:
	_, err = testee.SetTicketState(2, Closed)
	assert.Nil(t, err)

	// Assert that the ticketData has been updated:
	_, updatedTicket := testee.GetTicketById(testTicket.info.Id)
	assert.Equal(t, Closed, updatedTicket.info.State, "State should now be set to closed")

	// Assert that the stored file has been updated:
	storedTicket, err := readTicketFromFile(updatedTicket.filePath)
	assert.Equal(t, Closed, storedTicket.info.State, "The persisted state should be closed")
}

/*
	Setting the state of a non existing ticketData should return a error..
*/
func TestTicketManager_SetTicketState_TicketDoesNotExist(t *testing.T) {
	folderPath, rootPath, err := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	assert.Nil(t, err)

	// Load some prepared tickets:
	err = writeTestDataToFolder(folderPath)
	assert.Nil(t, err)

	testee := TicketManager{}
	testee.Initialize(folderPath)

	// Update the state:
	_, err = testee.SetTicketState(9999, Closed)
	assert.Equal(t, "ticketData does not exist", err.Error())
}

/*
	Write the test data to a folder.
*/
func writeTestDataToFolder(folderPath string) error {
	sampleData := []byte(firstTestTicket)
	// First ticketData: id 1
	sampleDataPath := path.Join(folderPath, "1.json")
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err := ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}

	// Second ticketData: id 2
	sampleData = []byte(secondTestTicket)
	sampleDataPath = path.Join(folderPath, "2.json")
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err = ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}

	// Third ticketData: id 3
	sampleData = []byte(thirdTestTicket)
	sampleDataPath = path.Join(folderPath, "3.json")
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err = ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}

	// Fourth ticketData: id 4
	sampleData = []byte(fourthTestTicket)
	sampleDataPath = path.Join(folderPath, "4.json")
	os.MkdirAll(filepath.Dir(sampleDataPath), 0644)
	err = ioutil.WriteFile(sampleDataPath, sampleData, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write test data file")
	}
	return nil
}

/*
	Prepare a temporary directory for tests.
*/
func prepareTempDirectory() (string, string, error) {
	// Creating a temp directory and remove it after the test:
	rootPath, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(rootPath)
	if err != nil {
		return "", "", err
	}
	// Create a path in the temp directory, the following path to a folder will not exist:
	notExistingFolderPath := path.Join(rootPath, "testDirectory")
	return notExistingFolderPath, rootPath, nil
}

/*
	Read a ticketData from a given file.
*/
func readTicketFromFile(filePath string) (*Ticket, error) {
	ticket, err := initializeFromFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not load ticker")
	}
	return ticket, nil
}

/*
	Asserting a ticketData.
*/
func assertTicket(t *testing.T, expected *Ticket, actual *Ticket) {
	// Id is set, so this can not be compared:
	// FilePath can not be compared
	// Creation- and LastModificationTime can not be compared

	// Compare ticketData info
	assert.Equal(t, expected.info.Title, actual.info.Title, "the title should be stored")
	assert.Equal(t, expected.info.Editor, actual.info.Editor, "the editor should be stored")
	assert.Equal(t, expected.info.HasEditor, actual.info.HasEditor, "the hasEditor field should be stored")
	assert.Equal(t, expected.info.Creator, actual.info.Creator, "the creator should be stored")

	// Compare messages
	assert.Equal(t, len(expected.messages), len(actual.Messages()))
	sort.Slice(expected.messages, func(i, j int) bool {
		return expected.messages[i].Id < expected.messages[j].Id
	})
	sort.Slice(actual.messages, func(i, j int) bool {
		return actual.messages[i].Id < actual.messages[j].Id
	})

	for i, expectedMessage := range expected.messages {
		actualMessage := actual.messages[i]
		// Creation time can not be compared
		assert.Equal(t, expectedMessage.Id, actualMessage.Id, "message id should be equal")
		assert.Equal(t, expectedMessage.CreatorMail, actualMessage.CreatorMail, "creator should be equal")
		assert.Equal(t, expectedMessage.OnlyInternal, actualMessage.OnlyInternal, " only internal should be equal")
		assert.Equal(t, expectedMessage.Content, actualMessage.Content, "content should be equal")
	}
}
func initializeWithTempTicketForExample(c *TicketManager) {
	folderPath, rootPath, _ := prepareTempDirectory()
	defer os.RemoveAll(rootPath)
	c.Initialize(folderPath)

	creator := Creator{Mail: "test@test.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	c.CreateNewTicket("Example_TicketTitle", creator, initialMessage)
}

func initializeWithTicketForExample(c *TicketManager) string {
	folderPath, rootPath, _ := prepareTempDirectory()
	c.Initialize(folderPath)

	creator := Creator{Mail: "test@test.de", FirstName: "Alex", LastName: "Wagner"}
	initialMessage := MessageEntry{Id: 1, CreatorMail: creator.Mail, Content: "This is a test", OnlyInternal: false}
	c.CreateNewTicket("Example_TicketTitle", creator, initialMessage)
	return rootPath
}
