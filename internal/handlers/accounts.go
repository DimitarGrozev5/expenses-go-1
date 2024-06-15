package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/accountsview"
	"github.com/go-chi/chi"
)

func (m *Repository) Accounts(w http.ResponseWriter, r *http.Request) {

	// Get all accounts
	accounts, err := m.DBClient.GetAccounts(r.Context(), &models.GetAccountsParams{OrderByPopularity: false})
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting accounts")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	// Get tags
	tags, err := m.DBClient.GetTags(r.Context(), nil)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting accounts")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	// Get user
	user, err := m.DBClient.GetUser(r.Context(), nil)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting accounts")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	// Get template data
	td := models.TemplateData{
		Title: "Accounts",
		Form: map[string]*forms.Form{
			"add-account":       forms.New(nil),
			"modify-free-funds": forms.New(nil),
		},
	}

	// Add forms for expenses
	for _, account := range accounts.Accounts {
		// Get form names
		moveUp := fmt.Sprintf("move-up-%d", account.ID)
		moveDown := fmt.Sprintf("move-down-%d", account.ID)
		delete := fmt.Sprintf("delete-%d", account.ID)

		// Add forms
		td.Form[moveUp] = forms.NewFromMap(map[string]string{
			"table_order": fmt.Sprintf("%d", account.TableOrder),
		})
		td.Form[moveDown] = forms.NewFromMap(map[string]string{
			"table_order": fmt.Sprintf("%d", account.TableOrder),
		})
		td.Form[delete] = forms.New(nil)
	}

	// Add default data
	m.AddDefaultData(&td, r)

	// Setup page data
	data := accountsview.AccountsData{
		TemplateData: td,
		Accounts:     accounts.Accounts,
		FreeFunds:    user.FreeFunds,
		Tags:         tags.Tags,
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

	// Get form and validate fields
	form := forms.New(r.PostForm)
	form.Required("name")
	form.MinLength("name", 4)

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
	name := form.Get("name")

	// Add expense to database
	_, err = m.DBClient.AddAccount(r.Context(), &models.AddAccountParams{Name: name})
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

func (m *Repository) PostModifyFreeFunds(w http.ResponseWriter, r *http.Request) {

	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// Get form and validate fields
	form := forms.New(r.PostForm)
	form.Required("amount", "to-account", "tag")
	form.IsFloat64("amount")
	form.IsInt("to-account")

	if !form.Valid() {

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"modify-free-funds": form,
		})

		// Redirect to expenses
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	// Get data
	amount, _ := strconv.ParseFloat(form.Get("amount"), 64)
	toAccount, _ := strconv.ParseInt(form.Get("to-account"), 10, 64)
	tag := form.Get("tag")

	// Add expense to database
	_, err = m.DBClient.ModifyFreeFunds(r.Context(), &models.ModifyFreeFundsParams{Amount: amount, ToAccountId: toAccount, TagName: tag})
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Failed to modify free funds")
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	// Add success message
	m.AddFlashMsg(r, "Free funds updated")
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
}

func (m *Repository) PostMoveAccount(direction int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		// Get account id from route param
		idParam := chi.URLParam(r, "accountId")
		id, err := strconv.ParseInt(idParam, 10, 32)
		if idParam == "" || err != nil {
			m.AddErrorMsg(r, "Invalid account")
			http.Redirect(w, r, "/accounts", http.StatusSeeOther)
			return
		}

		// Get form and validate fields
		form := forms.New(r.PostForm)
		form.Required("table_order")
		form.IsInt("table_order")

		if !form.Valid() {

			// Get table name
			directionWord := "up"
			if direction < 0 {
				directionWord = "down"
			}

			// Get form name
			name := fmt.Sprintf("move-%s-%d", directionWord, id)

			// Push form to session
			m.AddForms(r, map[string]*forms.Form{
				name: form,
			})

			// Redirect to expenses
			http.Redirect(w, r, "/accounts", http.StatusSeeOther)
			return
		}

		// Get data from form
		tableOrder, _ := strconv.ParseInt(form.Get("table_order"), 10, 64)

		// Construct account
		account := &models.GrpcAccount{
			ID:         id,
			TableOrder: tableOrder,
		}

		// Update account position
		_, err = m.DBClient.ReorderAccount(r.Context(), &models.ReorderAccountParams{Account: account, Direction: int64(direction)})
		if err != nil {
			m.App.ErrorLog.Println(err)
			m.AddErrorMsg(r, "Failed to move account")
			http.Redirect(w, r, "/accounts", http.StatusSeeOther)
			return
		}

		// Redirect
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
	}
}

func (m *Repository) PostDeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// Get account id from route param
	idParam := chi.URLParam(r, "accountId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if idParam == "" || err != nil {
		m.AddErrorMsg(r, "Invalid account")
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	// Delete account from database
	_, err = m.DBClient.DeleteAccount(r.Context(), &models.DeleteAccountParams{ID: id})
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Failed to delete account")
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	// Add success message
	m.AddFlashMsg(r, "Account deleted")
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
}
