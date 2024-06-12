package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/expensesview"
)

func (m *Repository) Expenses(w http.ResponseWriter, r *http.Request) {

	// Get db repo
	// repo, ok := m.GetDB(r)
	// if !ok {
	// 	m.App.ErrorLog.Println("Cannot get DB repo")
	// 	m.AddErrorMsg(r, "Please login to view expenses")
	// 	http.Redirect(w, r, "/logout", http.StatusSeeOther)
	// }

	// Get accounts
	accounts, err := m.DBClient.GetAccounts(r.Context(), &models.GetAccountsParams{OrderByPopularity: true})
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	// If there are no accounts redirect to accounts page and prompt user to create an account
	if len(accounts.Accounts) == 0 {
		m.AddWarningMsg(r, "You have to create a Payment Account before you can add Expenses")
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	// Get categories
	categories, err := m.DBClient.GetCategories(r.Context(), nil)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	// If there are no categories redirect to categories page and prompt user to create a category
	if len(categories.Categories) == 0 {
		m.AddWarningMsg(r, "You have to create a Budget Category before you can add Expenses")
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Get all expenses
	expenses, err := m.DBClient.GetExpenses(r.Context(), nil)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get all tags
	tags, err := m.DBClient.GetTags(r.Context(), nil)
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
			// "add-expense": forms.NewFromMap(map[string]string{
			// 	"amount":        "1",
			// 	"tags":          "asd",
			// 	"from_account":  "1",
			// 	"from_category": "1",
			// }),
		},
	}

	// Add forms for expenses
	for _, expense := range expenses.Expenses {
		// Get form names
		edit := fmt.Sprintf("edit-%d", expense.ID)
		delete := fmt.Sprintf("delete-%d", expense.ID)

		// Get expense tags as []string
		tags := make([]string, 0, len(expense.Tags))
		for _, tag := range expense.Tags {
			tags = append(tags, tag.Name)
		}

		date := expense.Date.AsTime()

		// Add forms
		td.Form[edit] = forms.NewFromMap(map[string]string{
			"amount":        fmt.Sprintf("%0.2f", expense.Amount),
			"tags":          strings.Join(tags, ","),
			"date":          fmt.Sprintf("%d-%02d-%02dT%02d:%02d", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute()),
			"from_account":  fmt.Sprintf("%d", expense.FromAccount.ID),
			"from_category": fmt.Sprintf("%d", expense.FromCategory.ID),
		})
		td.Form[delete] = forms.New(nil)
	}

	// Add default data
	m.AddDefaultData(&td, r)

	// Setup page data
	data := expensesview.ExpensesData{
		TemplateData: td,
		Expenses:     expenses.Expenses,
		Tags:         tags.Tags,
		Accounts:     accounts.Accounts,
		Categories:   categories.Categories,
	}

	// Render view
	data.View().Render(r.Context(), w)
}

func (m *Repository) PostNewExpense(w http.ResponseWriter, r *http.Request) {

	// // Parse form
	// err := r.ParseForm()
	// if err != nil {
	// 	log.Println(err)
	// }

	// // Get form and validate fields
	// form := forms.New(r.PostForm)
	// form.Required("amount", "tags", "from_account", "from_category", "date")
	// form.IsFloat64("amount")
	// form.MinLength("tags", 3)
	// form.IsDate("date", "2006-01-02T15:04")

	// if !form.Valid() {

	// 	// Push form to session
	// 	m.AddForms(r, map[string]*forms.Form{
	// 		"add-expense": form,
	// 	})

	// 	// Redirect to expenses
	// 	http.Redirect(w, r, "/expenses", http.StatusSeeOther)

	// 	return
	// }

	// // Get data
	// amount, _ := strconv.ParseFloat(form.Get("amount"), 64)
	// fromAccountId, _ := strconv.ParseInt(form.Get("from_account"), 10, 64)
	// fromCategoryId, _ := strconv.ParseInt(form.Get("from_category"), 10, 64)
	// date, _ := time.Parse("2006-01-02T15:04", form.Get("date"))

	// // Get tags
	// re := regexp.MustCompile(`,\s*`)

	// // Get db repo
	// repo, ok := m.GetDB(r)
	// if !ok {
	// 	m.App.ErrorLog.Println("Cannot get DB repo")
	// 	m.AddErrorMsg(r, "Please login to view expenses")
	// 	http.Redirect(w, r, "/logout", http.StatusSeeOther)
	// }

	// // Split tags field
	// tags := re.Split(form.Get("tags"), -1)

	// expense := models.Expense{
	// 	Amount:         amount,
	// 	Date:           date,
	// 	FromAccountId:  int(fromAccountId),
	// 	FromCategoryId: int(fromCategoryId),
	// }

	// // Add expense to database
	// err = repo.AddExpense(expense, tags)
	// if err != nil {
	// 	m.App.ErrorLog.Println(err)
	// 	m.AddErrorMsg(r, "Failed to add expense")
	// 	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	// 	return
	// }

	// // Add success message
	// m.AddFlashMsg(r, "Expense added")
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// Edit expense
func (m *Repository) PostEditExpense(w http.ResponseWriter, r *http.Request) {

	// // Parse form
	// err := r.ParseForm()
	// if err != nil {
	// 	log.Println(err)
	// }

	// // Get expense id from route param
	// idParam := chi.URLParam(r, "expenseId")
	// id, err := strconv.ParseInt(idParam, 10, 32)
	// if idParam == "" || err != nil {
	// 	m.AddErrorMsg(r, "Invalid expense")
	// 	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	// 	return
	// }

	// // Get form and validate fields
	// form := forms.New(r.PostForm)
	// form.Required("amount", "tags", "from_account", "from_category", "date")
	// form.IsFloat64("amount")
	// form.MinLength("tags", 3)
	// form.IsDate("date", "2006-01-02T15:04")

	// if !form.Valid() {

	// 	// Push form to session
	// 	m.AddForms(r, map[string]*forms.Form{
	// 		fmt.Sprintf("edit-%d", id): form,
	// 	})

	// 	// Redirect to expenses
	// 	http.Redirect(w, r, "/expenses", http.StatusSeeOther)

	// 	return
	// }

	// // Get data
	// amount, _ := strconv.ParseFloat(form.Get("amount"), 64)
	// fromAccountId, _ := strconv.ParseInt(form.Get("from_account"), 10, 64)
	// fromCategoryId, _ := strconv.ParseInt(form.Get("from_category"), 10, 64)
	// date, _ := time.Parse("2006-01-02T15:04", form.Get("date"))

	// // Get tags
	// re := regexp.MustCompile(`,\s*`)

	// // Get db repo
	// repo, ok := m.GetDB(r)
	// if !ok {
	// 	m.App.ErrorLog.Println("Cannot get DB repo")
	// 	m.AddErrorMsg(r, "Please login to view expenses")
	// 	http.Redirect(w, r, "/logout", http.StatusSeeOther)
	// }

	// // Split tags field
	// tags := re.Split(form.Get("tags"), -1)

	// expense := models.Expense{
	// 	ID:             int(id),
	// 	Amount:         amount,
	// 	Date:           date,
	// 	FromAccountId:  int(fromAccountId),
	// 	FromCategoryId: int(fromCategoryId),
	// }

	// // Add expense to database
	// err = repo.EditExpense(expense, tags)
	// if err != nil {
	// 	m.App.ErrorLog.Println(err)
	// 	m.AddErrorMsg(r, "Failed to edit expense")
	// 	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	// 	return
	// }

	// // Add success message
	// m.AddFlashMsg(r, "Expense updated")
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// Delete expense
func (m *Repository) PostDeleteExpense(w http.ResponseWriter, r *http.Request) {

	// // Parse form
	// err := r.ParseForm()
	// if err != nil {
	// 	log.Println(err)
	// }

	// // Get expense id from route param
	// idParam := chi.URLParam(r, "expenseId")
	// id, err := strconv.ParseInt(idParam, 10, 32)
	// if idParam == "" || err != nil {
	// 	m.AddErrorMsg(r, "Invalid expense")
	// 	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	// 	return
	// }

	// // Get DB repo
	// repo, ok := m.GetDB(r)
	// if !ok {
	// 	m.App.ErrorLog.Println("Failed to get DB repo")
	// 	m.AddErrorMsg(r, "Log in before adding expenses")
	// 	http.Redirect(w, r, "/logout", http.StatusSeeOther)
	// }

	// // Delete expense from database
	// err = repo.DeleteExpense(int(id))
	// if err != nil {
	// 	m.App.ErrorLog.Println(err)
	// 	m.AddErrorMsg(r, "Failed to delete expense")
	// 	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	// 	return
	// }

	// // Add success message
	// m.AddFlashMsg(r, "Expense deleted")
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
