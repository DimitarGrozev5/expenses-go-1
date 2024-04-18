package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/expensesview"
)

func (m *Repository) Expenses(w http.ResponseWriter, r *http.Request) {

	// Get db repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Cannot get DB repo")
		m.AddErrorMsg(r, "Please login to view expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get all expenses
	expenses, err := repo.GetExpenses()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get template data
	td := models.TemplateData{
		Title: "Expenses",
		Form: map[string]*forms.Form{
			"add-expense": forms.New(nil),
		},
	}

	// Add default data
	m.AddDefaultData(&td, r)

	// Setup page data
	data := expensesview.ExpensesData{
		TemplateData: td,
	}

	// Render view
	data.View().Render(r.Context(), w)
}

func (m *Repository) PostNewExpense(w http.ResponseWriter, r *http.Request) {

	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// Get form and validate fields
	form := forms.New(r.PostForm)
	form.Required("amount", "label", "date")
	form.IsFloat64("amount")
	form.MinLength("label", 3)
	form.IsDate("date", "2006-01-02T15:04")

	if !form.Valid() {

		// Get template data
		td := models.TemplateData{
			Title: "Expenses",
			Form: map[string]*forms.Form{
				"add-expense": form,
			},
		}

		// Add default data
		m.AddDefaultData(&td, r)

		// Setup page data
		data := expensesview.ExpensesData{
			TemplateData: td,
		}

		// Render view
		data.View().Render(r.Context(), w)
		return
	}

	// Get data
	amount, _ := strconv.ParseFloat(form.Get("amount"), 64)
	label := form.Get("label")
	date, _ := time.Parse("2006-01-02T15:04", form.Get("date"))

	expense := models.Expense{
		Amount: amount,
		Label:  label,
		Date:   date,
	}

	// Get DB repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Failed to get DB repo")
		m.AddErrorMsg(r, "Log in before adding expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Add expense to database
	err = repo.AddExpense(expense)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Failed to add expense")
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
		return
	}

	// Add success message
	m.AddFlashMsg(r, "Expense added")
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
