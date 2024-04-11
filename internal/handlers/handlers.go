package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/render"
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
		log.Println(err)
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
		form.IsFloat64("amount", r)
		form.MinLength("tag", 3, r)

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

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, r, "about.page.htm", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}

// var Static = http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
