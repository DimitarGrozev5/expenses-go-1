package handlers

import (
	"net/http"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/helpers"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/homeview"
)

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	// Check if authenticated and redirect to expenses page
	if helpers.IsAuthenticated(r) {
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	}

	// Get template data
	td := models.TemplateData{
		Title: "Home",
		Form: map[string]*forms.Form{
			"login": forms.New(nil),
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
