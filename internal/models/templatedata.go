package models

import "github.com/dimitargrozev5/expenses-go-1/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {

	// Deprecated
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}

	// CSRF Token
	CSRFToken string

	// Toast messages to flash on screen
	Flash   string
	Warning string
	Error   string

	// Page forms
	Form map[string]*forms.Form

	// Page title
	Title string

	// Page authentication status
	IsAuthenticated bool

	// Page url path
	CurrentURLPath string
}

// Check if dialog should be opened in templ dialog elements
func (d TemplateData) DialogOpened(dialog string) bool {
	return len(d.Form[dialog].Values) > 0
}
