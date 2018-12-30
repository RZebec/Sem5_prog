package templateManager

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
	A logger for tests.
*/
func getTestLogger() logging.Logger {
	return logging.ConsoleLogger{SetTimeStamp: false}
}

/*
	Tests if the Page is correctly rendered.
*/
func TestAccessDeniedPageRendering(t *testing.T) {
	rr := httptest.NewRecorder()

	mockConsoleLogger := getTestLogger()

	testee := TemplateManager{map[string]*template.Template{}}

	testee.LoadTemplates(mockConsoleLogger)
	err := testee.RenderTemplate(rr, "AccessDeniedPage", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := strings.TrimSpace(accessDeniedResultPage)
	result := strings.TrimSpace(rr.Body.String())
	if result != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			result, expected)
	}
}

/*
	Tests if the page returns an error if an page wasn´t able to be rendered.
*/
func TestErrorHandling(t *testing.T) {
	rr := httptest.NewRecorder()

	expectedError := NewError("Template doesn't exist")

	expectedText := "The template TestPage does not exist."

	testee := TemplateManager{map[string]*template.Template{}}

	err := testee.RenderTemplate(rr, "TestPage", nil)

	if err.Error() != expectedError.Error() {
		t.Errorf("error returned unexpected text: got %v want %v",
			err, expectedError)
	}

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	result := strings.TrimSpace(rr.Body.String())
	if result != expectedText {
		t.Errorf("handler returned unexpected body: got %v want %v",
			result, expectedText)
	}
}
