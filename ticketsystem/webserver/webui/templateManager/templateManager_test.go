package templateManager

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
	Tests if the Page is correctly rendered.
 */
func TestAccessDeniedPageRendering(t *testing.T) {
	rr := httptest.NewRecorder()

	LoadTemplates()
	err := RenderTemplate(rr, "AccessDeniedPage", nil)

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
	Mock structure for the Login Page Data.
*/
type mockLoginPageData struct {
	IsLoginFailed bool
}

/*
	Tests if the Page is correctly rendered for the given data.
 */
func TestLoginPageRendering(t *testing.T) {
	rr := httptest.NewRecorder()

	testData := mockLoginPageData{
		IsLoginFailed: false,
	}

	LoadTemplates()
	err := RenderTemplate(rr, "LoginPage", testData)

	if err != nil {
		t.Fatal(err)
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := strings.TrimSpace(loginResultPage)
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

	err := RenderTemplate(rr, "TestPage", nil)

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
