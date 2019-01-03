package templateManager

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"github.com/stretchr/testify/mock"
	"net/http"
)

/*
	A mocked template manager.
*/
type MockedTemplateManager struct {
	mock.Mock
}

/*
	Mocked function.
*/
func (m *MockedTemplateManager) LoadTemplates(logger logging.Logger) (err error) {
	args := m.Called(logger)
	return args.Error(0)
}

/*
	Mocked function.
*/
func (m *MockedTemplateManager) RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	args := m.Called(w, name, data)
	return args.Error(0)
}
