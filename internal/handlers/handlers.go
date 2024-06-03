package handlers

import (
	"net/http"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/helpers"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
	"github.com/justinas/nosurf"
)

// Repository used by the handlers
var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
	DB  map[string]repository.DatabaseRepo
}

// Creates a new repsoitory
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  make(map[string]repository.DatabaseRepo),
	}
}

// Sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Get DB connection
func (m *Repository) GetDB(r *http.Request) (repository.DatabaseRepo, bool) {
	db, ok := m.DB[m.App.Session.GetString(r.Context(), "user_key")]
	return db, ok
}

// Close all DB connections
func (m *Repository) CloseAllConnections() {
	for _, dbconn := range m.DB {
		dbconn.Close()
	}
}

// Add default data to template
func (m *Repository) AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Flash = m.App.Session.PopString(r.Context(), "flash")
	td.Error = m.App.Session.PopString(r.Context(), "error")
	td.Warning = m.App.Session.PopString(r.Context(), "warning")
	td.IsAuthenticated = helpers.IsAuthenticated(r)
	td.CurrentURLPath = r.URL.Path

	if m.App.Session.Exists(r.Context(), "forms") {
		storedForms := m.App.Session.Pop(r.Context(), "forms").(map[string]*forms.Form)
		for formName, form := range storedForms {
			td.Form[formName] = form
		}
	}

	return td
}

// Add flash message to session
func (m *Repository) AddFlashMsg(r *http.Request, msg string) {
	m.App.Session.Put(r.Context(), "flash", msg)
}

// Add warning message to session
func (m *Repository) AddWarningMsg(r *http.Request, msg string) {
	m.App.Session.Put(r.Context(), "warning", msg)
}

// Add error message to session
func (m *Repository) AddErrorMsg(r *http.Request, msg string) {
	m.App.Session.Put(r.Context(), "error", msg)
}

// Add form to session
func (m *Repository) AddForms(r *http.Request, forms map[string]*forms.Form) {
	m.App.Session.Put(r.Context(), "forms", forms)
}

func (m *Repository) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}
