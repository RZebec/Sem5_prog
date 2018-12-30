package templateManager

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/templateManager/pages"
	"fmt"
	"html/template"
	"net/http"
)

/*
	Inspiration from source: https://hackernoon.com/golang-template-2-template-composition-and-how-to-organize-template-files-4cb40bcdf8f6
*/

type TemplateContext interface {
	// Loads all available templates from their corresponding strings in the template map.
	LoadTemplates(logger logging.Logger) (err error)
	// Renders the needed template with the given name and the needed page data.
	RenderTemplate(w http.ResponseWriter, name string, data interface{}) error
}

var bufpool *helpers.BufferPool

// Create a buffer pool
func init() {
	bufpool = helpers.NewBufferPool(64)
}

type TemplateManager struct {
	Templates map[string]*template.Template
}

/*
	Struct for the template error.
*/
type TemplateError struct {
	message string
}

/*
	Function that returns the error message.
*/
func (e *TemplateError) Error() string {
	return e.message
}

/*
	Function for defining an template error.
*/
func NewError(text string) error {
	return &TemplateError{text}
}

/*
	Loads all available templates from their corresponding strings in the template map.
*/
func (t *TemplateManager) LoadTemplates(logger logging.Logger) (err error) {

	if t.Templates == nil {
		t.Templates = make(map[string]*template.Template)
	}

	baseTemplate := template.New("Base")

	baseTemplate, err = baseTemplate.Parse(pages.Base)

	if err != nil {
		return err
	}

	t.addTemplate(pages.IndexPage, "IndexPage", baseTemplate, logger)
	t.addTemplate(pages.RegisterPage, "RegisterPage", baseTemplate, logger)
	t.addTemplate(pages.LoginPage, "LoginPage", baseTemplate, logger)
	t.addTemplate(pages.AdminPage, "AdminPage", baseTemplate, logger)
	t.addTemplate(pages.AccessDeniedPage, "AccessDeniedPage", baseTemplate, logger)
	t.addTemplate(pages.TicketExplorerPage, "TicketExplorerPage", baseTemplate, logger)
	t.addTemplate(pages.TicketViewPage, "TicketViewPage", baseTemplate, logger)

	return nil
}

/*
	Helper function.
	Adds a template to the template map with the corresponding name and template string.
*/
func (t *TemplateManager) addTemplate(templateString string, templateName string, baseTemplate *template.Template, logger logging.Logger) {
	var err error

	t.Templates[templateName], err = baseTemplate.Clone()

	if err != nil {
		logger.LogError("Template", err)
	}

	t.Templates[templateName].New(templateName)

	t.Templates[templateName], err = t.Templates[templateName].Parse(templateString)

	if err != nil {
		logger.LogError("Template", err)
	}
}

/*
	Renders the needed template with the given name and the needed page data.
*/
func (t *TemplateManager) RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := t.Templates[name]

	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist.", name),
			http.StatusInternalServerError)
		err := NewError("Template doesn't exist")
		return err
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.Execute(buf, data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}
