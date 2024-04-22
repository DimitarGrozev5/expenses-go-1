package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/expensesview"
	"github.com/go-chi/chi"
)

func (m *Repository) Expenses(w http.ResponseWriter, r *http.Request) {

	// Get db repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Cannot get DB repo")
		m.AddErrorMsg(r, "Please login to view expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get all accounts
	accountsCount, err := repo.GetAccountsCount()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	// If there are not accounts redirect to accounts page and prompt user to create an account
	if accountsCount == 0 {
		m.AddWarningMsg(r, "You have to create a Payment Account before you can add Expenses")
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	// Get all expenses
	expenses, err := repo.GetExpenses()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get all tags
	tags, err := repo.GetTags()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting data")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get template data
	td := models.TemplateData{
		Title: "Expenses",
		Form: map[string]*forms.Form{
			"add-expense": forms.New(nil),
		},
	}

	// Add forms for expenses
	for _, expense := range expenses {
		// Get form names
		edit := fmt.Sprintf("edit-%d", expense.ID)
		delete := fmt.Sprintf("delete-%d", expense.ID)

		// Get expense tags as []string
		tags := make([]string, 0, len(expense.Tags))
		for _, tag := range expense.Tags {
			tags = append(tags, tag.Name)
		}

		// Add forms
		td.Form[edit] = forms.NewFromMap(map[string]string{
			"amount": fmt.Sprintf("%0.2f", expense.Amount),
			"tags":   strings.Join(tags, ","),
			"date":   fmt.Sprintf("%d-%02d-%02dT%02d:%02d", expense.Date.Year(), expense.Date.Month(), expense.Date.Day(), expense.Date.Hour(), expense.Date.Minute()),
		})
		td.Form[delete] = forms.New(nil)
	}

	// Add default data
	m.AddDefaultData(&td, r)

	// Setup page data
	data := expensesview.ExpensesData{
		TemplateData: td,
		Expenses:     expenses,
		Tags:         tags,
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
	form.Required("amount", "tags", "date")
	form.IsFloat64("amount")
	form.MinLength("tags", 3)
	form.IsDate("date", "2006-01-02T15:04")

	if !form.Valid() {

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"add-expense": form,
		})

		// Redirect to expenses
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)

		return
	}

	// Get data
	amount, _ := strconv.ParseFloat(form.Get("amount"), 64)
	date, _ := time.Parse("2006-01-02T15:04", form.Get("date"))

	// Get tags
	re := regexp.MustCompile(`,\s*`)

	// Get db repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Cannot get DB repo")
		m.AddErrorMsg(r, "Please login to view expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get all tags
	allTags, err := repo.GetTags()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting data")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Split tags field
	split := re.Split(form.Get("tags"), -1)

	// Get tags
	tags := make([]models.Tag, 0, len(split))
	for _, tagName := range split {

		// If tag is in all tags
		tagSet := false
		for _, tag := range allTags {
			if tag.Name == tagName {
				tags = append(tags, tag)
				tagSet = true
				break
			}
		}

		// If tag is not in all tags
		if !tagSet {
			tags = append(tags, models.Tag{
				ID:         -1,
				Name:       tagName,
				UsageCount: 1,
				LastUsed:   time.Now(),
			})
		}
	}

	expense := models.Expense{
		Amount: amount,
		Tags:   tags,
		Date:   date,
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

// Edit expense
func (m *Repository) PostEditExpense(w http.ResponseWriter, r *http.Request) {

	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// Get expense id from route param
	idParam := chi.URLParam(r, "expenseId")
	id, err := strconv.ParseInt(idParam, 0, 32)
	if idParam == "" || err != nil {
		m.AddErrorMsg(r, "Invalid expense")
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
		return
	}

	// Get form and validate fields
	form := forms.New(r.PostForm)
	form.Required("amount", "tags", "date")
	form.IsFloat64("amount")
	form.MinLength("tags", 3)
	form.IsDate("date", "2006-01-02T15:04")

	if !form.Valid() {

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"add-expense": form,
		})

		// Redirect to expenses
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)

		return
	}

	// Get data
	amount, _ := strconv.ParseFloat(form.Get("amount"), 64)
	date, _ := time.Parse("2006-01-02T15:04", form.Get("date"))

	// Get tags
	re := regexp.MustCompile(`,\s*`)

	// Get db repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Cannot get DB repo")
		m.AddErrorMsg(r, "Please login to view expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get all tags
	allTags, err := repo.GetTags()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting data")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Split tags field
	split := re.Split(form.Get("tags"), -1)

	// Get tags
	tags := make([]models.Tag, 0, len(split))
	for _, tagName := range split {

		// If tag is in all tags
		tagSet := false
		for _, tag := range allTags {
			if tag.Name == tagName {
				tags = append(tags, tag)
				tagSet = true
				break
			}
		}

		// If tag is not in all tags
		if !tagSet {
			tags = append(tags, models.Tag{
				ID:         -1,
				Name:       tagName,
				UsageCount: 1,
				LastUsed:   time.Now(),
			})
		}
	}

	expense := models.Expense{
		ID:     int(id),
		Amount: amount,
		Tags:   tags,
		Date:   date,
	}

	// Add expense to database
	err = repo.EditExpense(expense)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Failed to edit expense")
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
		return
	}

	// Add success message
	m.AddFlashMsg(r, "Expense updated")
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// Delete expense
func (m *Repository) PostDeleteExpense(w http.ResponseWriter, r *http.Request) {

	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// Get expense id from route param
	idParam := chi.URLParam(r, "expenseId")
	id, err := strconv.ParseInt(idParam, 0, 32)
	if idParam == "" || err != nil {
		m.AddErrorMsg(r, "Invalid expense")
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
		return
	}

	// Get DB repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Failed to get DB repo")
		m.AddErrorMsg(r, "Log in before adding expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Add expense to database
	err = repo.DeleteExpense(int(id))
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Failed to delete expense")
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
		return
	}

	// Add success message
	m.AddFlashMsg(r, "Expense deleted")
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
