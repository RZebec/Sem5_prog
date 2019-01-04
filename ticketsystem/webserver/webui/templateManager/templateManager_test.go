package templateManager

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/testhelpers"
	"github.com/stretchr/testify/assert"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
	Loading all templates should be possible.
*/
func TestTemplateManager_LoadTemplates(t *testing.T) {
	testee := TemplateManager{}

	err := testee.LoadTemplates(testhelpers.GetTestLogger())

	assert.Nil(t, err, "All templates should be loaded without error")
	assert.Equal(t, 7, len(testee.Templates), "All templates should be loaded")
}

/*
	Rendering a template should replace the placeholders with the data.
*/
func TestTemplateManager_RenderTemplate(t *testing.T) {
	testee := TemplateManager{}
	baseTemplate := template.New("Base")
	err := testee.addTemplate(testPage, "TestPage", baseTemplate, testhelpers.GetTestLogger())

	assert.Nil(t, err, "Template should be added without problem")

	testData := testPageData{Title: "TestTitleToInsert"}
	rr := httptest.NewRecorder()

	// Execute the test:
	err = testee.RenderTemplate(rr, "TestPage", testData)

	// Assert the result
	assert.Nil(t, err, "Template should be added without problem")
	resultContent := rr.Result()

	defer resultContent.Body.Close()
	body, err := ioutil.ReadAll(resultContent.Body)
	assert.Contains(t, string(body), "TestTitleToInsert", "Data should be set")
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
