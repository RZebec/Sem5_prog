package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/user"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

/*
	A Valid Http Request for logged In user.
*/
func TestTicketCreateHandler_ServeHTTP_ValidRequest(t *testing.T) {
	mail := "test@test.com"
	title := "Test Ticket"
	message := "This is a Test Ticket Message"
	internal := "false"

	req, err := http.NewRequest("POST", "/create_ticket", nil)
	req.Form = url.Values{}
	req.Form.Add("mail", mail)
	req.Form.Add("title", title)
	req.Form.Add("message", message)
	req.Form.Add("internal", internal)

	testUser := user.User{Mail: mail, UserId: 5, FirstName: "Max", LastName: "Muller", Role: user.RegisteredUser, State: user.Active}
	testMessage := ticket.MessageEntry{Id: 0, CreatorMail: mail, Content: message, OnlyInternal: false}

	testTicketCreator := ticket.Creator{Mail: mail, FirstName: "Max", LastName: "Muller"}

	testTicketInfo := ticket.TicketInfo{Id: 26, Title: title, Editor: testUser, HasEditor: true, Creator: testTicketCreator, CreationTime: time.Now(), LastModificationTime: time.Now(), State: ticket.Processing}

	ticket := ticket.CreateTestTicket(testTicketInfo, []ticket.MessageEntry{testMessage})

	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTicketContext := new(mockedForTests.MockedTicketContext)
	mockedTicketContext.On("CreateNewTicketForInternalUser", title, testUser, testMessage).Return(ticket, nil)
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedUserContext.On("GetUserForEmail", mail).Return(true, testUser.UserId)
	mockedUserContext.On("GetUserById", testUser.UserId).Return(true, testUser)

	testee := TicketCreateHandler{TicketContext: mockedTicketContext, UserContext: mockedUserContext, Logger: testLogger}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, 302, resp.StatusCode, "Should return status code 302")
	assert.Equal(t, "/ticket/26", newLocation, "Should be redirected to \"/ticket/26\"")

	mockedTicketContext.AssertExpectations(t)
	mockedUserContext.AssertExpectations(t)

}

