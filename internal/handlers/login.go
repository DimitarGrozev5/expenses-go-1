package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository/dbrepo"
)

// Handle posting to login
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	// Renew session token
	_ = m.App.Session.RenewToken(r.Context())

	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// Get form and validate fields
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"login": form,
		})

		// Redirect to home
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	// Get user data
	uEmail := form.Get("email")
	uPassword := form.Get("password")

	// Check if user DB exists
	_, err = os.Stat(dbrepo.GetUserDBPath(m.App.DBPath, uEmail))
	if errors.Is(err, os.ErrNotExist) {

		// Write to error log
		m.App.ErrorLog.Println(err)

		// Reset password in form
		form.Set("password", "")

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"login": form,
		})

		// Add error message
		m.AddErrorMsg(r, "Invalid login credentials")

		// Redirect to home
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	// Create user connection
	dbconn, err := driver.ConnectSQL(dbrepo.GetUserDBPath(m.App.DBPath, uEmail))
	if err != nil {

		// Write to error log
		m.App.ErrorLog.Println(err)

		// Reset password in form
		form.Set("password", "")

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"login": form,
		})

		// Add error message
		m.AddErrorMsg(r, "Server error")

		// Redirect to home
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	// Get db repo
	repo := dbrepo.NewSqliteRepo(m.App, uEmail, dbconn.SQL)

	// Authenticate user
	_, _, err = repo.Authenticate(uEmail, uPassword)
	if err != nil {

		// Write to error log
		m.App.ErrorLog.Println(err)

		// Reset password in form
		form.Set("password", "")

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"login": form,
		})

		// Add error message
		m.AddErrorMsg(r, "Invalid login credentials")

		// Redirect to home
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	// Get user key
	key := dbrepo.GetUserKey(uEmail)

	// Add connection to repo
	m.DB[key] = repo

	// Store user key in session
	m.App.Session.Put(r.Context(), "user_key", key)

	// Flash message to user
	m.AddFlashMsg(r, "Logged in successfully")

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
