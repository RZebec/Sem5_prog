package templateManager

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/webui/helpers"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var bufpool *helpers.BufferPool

// create a buffer pool
func init() {
	bufpool = helpers.NewBufferPool(64)
	log.Println("buffer allocation successful")
}

var templates map[string]*template.Template

type TemplateError struct {
	s string
}

func (e *TemplateError) Error() string {
	return e.s
}

func NewError(text string) error {
	return &TemplateError{text}
}

func LoadTemplates() (err error) {

	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	baseTemplate := template.New("Base")

	baseTemplate, err = baseTemplate.Parse(base)

	if err != nil {
		// TODO: Handle error
		fmt.Print(err)
	}

	addTemplate(loginPage, "LoginPage", baseTemplate)

	return nil
}

/*
	Parses multiple template Strings
	Source: https://stackoverflow.com/questions/41856021/how-to-parse-multiple-strings-into-a-template-with-go
 */
func parseTemplates(templs ...string) (t *template.Template, err error) {
	t = template.New("_all")

	for i, templ := range templs {
		if _, err = t.New(fmt.Sprint("_", i)).Parse(templ); err != nil {
			return
		}
	}

	return
}

func addTemplate(templateString string, templateName string, baseTemplate *template.Template) {
	var err error

	templates[templateName], err = baseTemplate.Clone()

	if err != nil {
		// TODO: Handle error
		fmt.Print(err)
	}

	templates[templateName].New(templateName)

	templates[templateName], err = templates[templateName].Parse(templateString)

	if err != nil {
		// TODO: Handle error
		fmt.Print(err)
	}
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]

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
		// TODO: Handle error
		fmt.Print(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}