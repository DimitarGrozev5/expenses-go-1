package handlers

import (
	"net/http"
)

// Handle posting to login
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user key from session
	key := m.App.Session.GetString(r.Context(), "user_key")

	// Destroy session and renew session token
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	// Close connection and delete from repo
	m.DB[key].Close()
	delete(m.DB, key)

	// Flash message to user
	m.App.Session.Put(r.Context(), "flash", "Logged out successfully")

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
