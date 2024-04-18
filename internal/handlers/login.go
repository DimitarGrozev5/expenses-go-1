package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository/dbrepo"
	"github.com/dimitargrozev5/expenses-go-1/views/homeview"
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

		// Get template data
		td := models.TemplateData{
			Title: "Home",
			Form: map[string]*forms.Form{
				"login": form,
			},
		}

		// Add default data
		m.AddDefaultData(&td, r)

		// Setup page data
		data := homeview.HomeData{
			TemplateData: td,
		}

		// Render view
		data.View().Render(r.Context(), w)
	}

	// Get user data
	uEmail := form.Get("email")
	uPassword := form.Get("password")

	// Check if user DB exists
	_, err = os.Stat(dbrepo.GetUserDBPath(m.App.DBPath, uEmail))
	if errors.Is(err, os.ErrNotExist) {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Create user connection
	dbconn, err := driver.ConnectSQL(dbrepo.GetUserDBPath(m.App.DBPath, uEmail))
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Server error")
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	// Get db repo
	repo := dbrepo.NewSqliteRepo(m.App, uEmail, dbconn.SQL)

	// Authenticate user
	_, _, err = repo.Authenticate(uEmail, uPassword)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
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
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
