package handlers

import (
	"fmt"
	"net/http"

	"github.com/dimitargrozev5/expenses-go-1/pkg/config"
	"github.com/dimitargrozev5/expenses-go-1/pkg/models"
	"github.com/dimitargrozev5/expenses-go-1/pkg/render"
)

// Repository used by the handlers
var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
}

// Creates a new repsoitory
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// Sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr

	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "home.page.htm", &models.TemplateData{})
}

func (m *Repository) Expenses(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "expenses.page.htm", &models.TemplateData{})
}

func (m *Repository) PostExpenses(w http.ResponseWriter, r *http.Request) {
	action := r.Form.Get("action")

	w.Write([]byte(fmt.Sprintf("Post Expenses. Type: %s", action)))
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, r, "about.page.htm", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}

// var Static = http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
