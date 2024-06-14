package handlers

import (
	"net/http"
)

// Handle posting to login
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {

	// Destroy session and renew session token
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	// Flash message to user
	m.App.Session.Put(r.Context(), "flash", "Logged out successfully")

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
