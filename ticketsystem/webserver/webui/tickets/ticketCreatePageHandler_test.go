package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mockedForTests"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/userData"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/wrappers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	A Valid Http Request for logged In userData.
*/
func TestTicketCreatePageHandler_ServeHTTP_ValidRequest(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := TicketCreatePageHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	testUser := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}

	req, err := http.NewRequest("GET", "/ticket_create", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")

	data := ticketCreatePageData{
		UserName:       testUser.Mail,
		IsUserLoggedIn: true,
		FirstName:      "Dieter",
		LastName:       "Dietrich",
	}

	data.UserIsAdmin = false
	data.UserIsAuthenticated = true
	data.Active = "ticket_create"

	mockedUserContext.On("GetUserById", 5).Return(true, testUser)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketCreatePage", data).Return(nil)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return status code 200")
	assert.Equal(t, "", newLocation, "Should not be redirected")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

/*
	Only GET request should be possible.
*/
func TestTicketCreatePageHandler_ServeHTTP_WrongRequestMethod(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := TicketCreatePageHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("POST", "/ticket_create", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Should return status code 405")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

/*
	A error from rendering the template should result in a 500.
*/
func TestTicketCreatePageHandler_ServeHTTP_ContextError_RenderError(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := TicketCreatePageHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	testUser := userData.User{Mail: "Test2@Test.de", UserId: 5, FirstName: "Dieter", LastName: "Dietrich", Role: userData.RegisteredUser, State: userData.Active}

	req, err := http.NewRequest("GET", "/ticket_create", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")

	data := ticketCreatePageData{
		UserName:       testUser.Mail,
		IsUserLoggedIn: true,
		FirstName:      "Dieter",
		LastName:       "Dietrich",
	}

	data.UserIsAdmin = false
	data.UserIsAuthenticated = true
	data.Active = "ticket_create"

	mockedUserContext.On("GetUserById", 5).Return(true, testUser)
	mockedTemplateManager.On("RenderTemplate", mock.Anything, "TicketCreatePage", data).Return(errors.New("TestError"))

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/", newLocation, "Should be redirected to \"/\"")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}

/*
	If the logged in userData doesn´t exist it should result in a 500.
*/
func TestTicketCreatePageHandler_ServeHTTP_UserDoesNotExist(t *testing.T) {
	mockedUserContext := new(mockedForTests.MockedUserContext)
	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := TicketCreatePageHandler{UserContext: mockedUserContext, TemplateManager: mockedTemplateManager,
		Logger: testhelpers.GetTestLogger()}

	req, err := http.NewRequest("GET", "/ticket_create", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute the test:
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), true, false, 5, "")

	mockedUserContext.On("GetUserById", 5).Return(false, *new(userData.User))

	handler.ServeHTTP(rr, req.WithContext(ctx))

	resp := rr.Result()

	newLocation := resp.Header.Get("location")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Should return status code 500")
	assert.Equal(t, "/", newLocation, "Should be redirected to \"/\"")

	mockedUserContext.AssertExpectations(t)
	mockedTemplateManager.AssertExpectations(t)
}
