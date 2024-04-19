package handlers

import (
	"net/http"
	"strconv"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/helpers"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/render"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
	"github.com/dimitargrozev5/expenses-go-1/views"
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

func (m *Repository) Home1(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.htm", &models.TemplateData{})
	views.Test().Render(r.Context(), w)
	// layout.MainLayout("").Render(r.Context(), w)
}

func (m *Repository) Expenses1(w http.ResponseWriter, r *http.Request) {
	var emptyExpense models.NewExpense
	data := make(map[string]interface{})
	data["newExpense"] = emptyExpense

	render.RenderTemplate(w, r, "expenses.page.htm", &models.TemplateData{
		Form: map[string]*forms.Form{
			"create": forms.New(nil),
			"update": forms.New(nil),
			"delete": forms.New(nil),
		},
	})
}

func (m *Repository) PostExpenses(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	action := r.Form.Get("action")

	switch action {
	case "create":
		amount, _ := strconv.ParseFloat(r.Form.Get("amount"), 64)
		expense := models.NewExpense{
			Amount: amount,
			Tag:    r.Form.Get("tag"),
			Date:   r.Form.Get("date"),
		}

		form := forms.New(r.PostForm)

		form.Required("amount", "tag", "date")
		form.IsFloat64("amount")
		form.MinLength("tag", 3)

		if !form.Valid() {
			data := make(map[string]interface{})
			data["newExpense"] = expense
			data["createOpen"] = "true"

			render.RenderTemplate(w, r, "expenses.page.htm", &models.TemplateData{
				Form: map[string]*forms.Form{
					"create": form,
					"update": forms.New(nil),
					"delete": forms.New(nil),
				},
				Data: data,
			})

			return
		}

		m.App.Session.Put(r.Context(), "flash", "Expense added")
	}

	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	render.RenderTemplate(w, r, "about.page.htm", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}
