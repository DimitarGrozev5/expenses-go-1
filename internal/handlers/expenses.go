package handlers

import (
	"net/http"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/expensesview"
)

func (m *Repository) Expenses(w http.ResponseWriter, r *http.Request) {

	// Get template data
	td := models.TemplateData{
		Title: "Expenses",
		// Form: map[string]*forms.Form{
		// 	"login": forms.New(nil),
		// },
	}

	// Add default data
	m.AddDefaultData(&td, r)

	// Setup page data
	data := expensesview.ExpensesData{
		TemplateData: td,
	}

	// Render view
	data.View().Render(r.Context(), w)
}
