package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/expensesview"
)

func (m *Repository) Expenses(w http.ResponseWriter, r *http.Request) {

	// Get template data
	td := models.TemplateData{
		Title: "Expenses",
		Form: map[string]*forms.Form{
			"add-expense": forms.New(nil),
		},
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

func (m *Repository) PostNewExpense(w http.ResponseWriter, r *http.Request) {

	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// Get form and validate fields
	form := forms.New(r.PostForm)
	form.Required("amount", "label", "date")
	form.IsFloat64("amount")
	form.MinLength("label", 3)
	form.IsDate("date", "2006-01-02T15:04")

	if !form.Valid() {

		// Get template data
		td := models.TemplateData{
			Title: "Expenses",
			Form: map[string]*forms.Form{
				"add-expense": form,
			},
		}

		// Add default data
		m.AddDefaultData(&td, r)

		// Setup page data
		data := expensesview.ExpensesData{
			TemplateData: td,
		}

		// Render view
		data.View().Render(r.Context(), w)
		return
	}

	// Get data
	amount, _ := strconv.ParseFloat(form.Get("amount"), 64)
	label := form.Get("label")
	date, _ := time.Parse("2006-01-02T15:04", form.Get("date"))

	fmt.Println(amount, label, date)
}
