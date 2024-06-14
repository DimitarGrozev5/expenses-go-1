package handlers

import (
	"log"
	"net/http"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
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

	// Authenticate user
	result, err := m.DBClient.Authenticate(r.Context(), &models.LoginCredentials{Email: uEmail, Password: uPassword})
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

	// Store user token in session
	m.App.Session.Put(r.Context(), "user_token", result.Token)

	// Flash message to user
	m.AddFlashMsg(r, "Logged in successfully")

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
