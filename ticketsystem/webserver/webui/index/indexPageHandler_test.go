// 5894619, 6720876, 9793350
package index

import (
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
	A valid request should be possible.
*/
func TestIndexPageHandler_ServeHTTP_ValidRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	rr := httptest.NewRecorder()

	testee := PageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, "Status code 200 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	A error from rendering the template should result in a 500.
*/
func TestIndexPageHandler_ServeHTTP_RenderTemplateError_500Returned(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	mockedTemplateManager.On("RenderTemplate", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("TestError"))

	rr := httptest.NewRecorder()

	testee := PageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code, "Status code 500 should be returned")

	mockedTemplateManager.AssertExpectations(t)
}

/*
	Only GET Methods should be allowed.
*/
func TestIndexPageHandler_ServeHTTP_WrongRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	testLogger := testhelpers.GetTestLogger()

	mockedTemplateManager := new(templateManager.MockedTemplateManager)

	testee := PageHandler{Logger: testLogger, TemplateManager: mockedTemplateManager}

	handler := http.HandlerFunc(testee.ServeHTTP)
	ctx := wrappers.NewContextWithAuthenticationInfo(req.Context(), false, false, -1, "")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Status code 405 should be returned")
}
