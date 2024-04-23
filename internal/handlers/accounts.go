package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/accountsview"
)

func (m *Repository) Accounts(w http.ResponseWriter, r *http.Request) {
	// Get db repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Cannot get DB repo")
		m.AddErrorMsg(r, "Please login to view expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get all accounts
	accounts, err := repo.GetAccounts()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting accounts")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get template data
	td := models.TemplateData{
		Title: "Accounts",
		Form: map[string]*forms.Form{
			"add-account": forms.New(nil),
		},
	}

	// // Add forms for expenses
	// for _, expense := range expenses {
	// 	// Get form names
	// 	edit := fmt.Sprintf("edit-%d", expense.ID)
	// 	delete := fmt.Sprintf("delete-%d", expense.ID)

	// 	// Get expense tags as []string
	// 	tags := make([]string, 0, len(expense.Tags))
	// 	for _, tag := range expense.Tags {
	// 		tags = append(tags, tag.Name)
	// 	}

	// 	// Add forms
	// 	td.Form[edit] = forms.NewFromMap(map[string]string{
	// 		"amount": fmt.Sprintf("%0.2f", expense.Amount),
	// 		"tags":   strings.Join(tags, ","),
	// 		"date":   fmt.Sprintf("%d-%02d-%02dT%02d:%02d", expense.Date.Year(), expense.Date.Month(), expense.Date.Day(), expense.Date.Hour(), expense.Date.Minute()),
	// 	})
	// 	td.Form[delete] = forms.New(nil)
	// }

	// Add default data
	m.AddDefaultData(&td, r)

	// Setup page data
	data := accountsview.AccountsData{
		TemplateData: td,
		Accounts:     accounts,
	}

	// Render view
	data.View().Render(r.Context(), w)
}

func (m *Repository) PostNewAccount(w http.ResponseWriter, r *http.Request) {

	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// Get db repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Cannot get DB repo")
		m.AddErrorMsg(r, "Please login to view expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get form and validate fields
	form := forms.New(r.PostForm)
	form.Required("name", "initial_amount")
	form.IsFloat64("initial_amount")
	form.MinLength("name", 3)

	if !form.Valid() {

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"add-account": form,
		})

		// Redirect to expenses
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)

		return
	}

	// Get data
	initialAmount, _ := strconv.ParseFloat(form.Get("initial_amount"), 64)
	name := form.Get("name")

	// Get Account object
	account := models.Account{
		Name:          name,
		InitialAmount: initialAmount,
	}

	// Add expense to database
	err = repo.AddAccount(account)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Failed to add account")
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	// Add success message
	m.AddFlashMsg(r, "Account added")
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
}
