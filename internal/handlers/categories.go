package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/categoriesview"
	"github.com/go-chi/chi"
)

func (m *Repository) Categories(w http.ResponseWriter, r *http.Request) {
	// Get db repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Cannot get DB repo")
		m.AddErrorMsg(r, "Please login to view expenses")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get all categories
	categories, err := repo.GetCategoriesOverview()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting categories")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get time periods
	periods, err := repo.GetTimePeriods()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting categories")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Get template data
	td := models.TemplateData{
		Title: "Categories",
		Form: map[string]*forms.Form{
			"add-category":     forms.New(nil),
			"reset-categories": forms.New(nil),
		},
	}

	// Add forms for expenses
	for _, category := range categories {
		// Get form names
		moveUp := fmt.Sprintf("move-up-%d", category.ID)
		moveDown := fmt.Sprintf("move-down-%d", category.ID)
		delete := fmt.Sprintf("delete-%d", category.ID)

		// Add forms
		td.Form[moveUp] = forms.NewFromMap(map[string]string{
			"table_order": fmt.Sprintf("%d", category.TableOrder),
		})
		td.Form[moveDown] = forms.NewFromMap(map[string]string{
			"table_order": fmt.Sprintf("%d", category.TableOrder),
		})
		td.Form[delete] = forms.New(nil)
	}

	// Add default data
	m.AddDefaultData(&td, r)

	// Setup page data
	data := categoriesview.CategoriesData{
		TemplateData: td,
		Categories:   categories,
		TimePeriods:  periods,
	}

	// Render view
	data.View().Render(r.Context(), w)
}

func (m *Repository) PostNewCategory(w http.ResponseWriter, r *http.Request) {

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
	form.Required("name", "budget_input", "spending_limit", "input_interval", "input_period")
	form.MinLength("name", 4)

	form.IsFloat64("budget_input")
	form.Min("budget_input", 0)

	form.IsFloat64("spending_limit")
	form.Min("spending_limit", 0)

	form.IsInt("input_interval")
	form.Min("input_interval", 1)

	form.IsInt("input_period")

	if !form.Valid() {

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"add-category": form,
		})

		// Redirect to categories
		http.Redirect(w, r, "/categories", http.StatusSeeOther)

		return
	}

	// Get data
	name := form.Get("name")
	budgetInput, _ := strconv.ParseFloat(form.Get("budget_input"), 64)
	spendingLimit, _ := strconv.ParseFloat(form.Get("spending_limit"), 64)
	inputInterval, _ := strconv.ParseInt(form.Get("input_interval"), 10, 64)
	inputPeriod, _ := strconv.ParseInt(form.Get("input_period"), 10, 64)

	// Add category to database
	err = repo.AddCategory(name, budgetInput, spendingLimit, int(inputInterval), int(inputPeriod))
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Failed to add category")
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Add success message
	m.AddFlashMsg(r, "Category added")
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (m *Repository) PostMoveCategory(direction int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		// Get category id from route param
		idParam := chi.URLParam(r, "categoryId")
		id, err := strconv.ParseInt(idParam, 10, 32)
		if idParam == "" || err != nil {
			m.AddErrorMsg(r, "Invalid category")
			http.Redirect(w, r, "/categories", http.StatusSeeOther)
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
			http.Redirect(w, r, "/categories", http.StatusSeeOther)

			return
		}

		// Get db repo
		repo, ok := m.GetDB(r)
		if !ok {
			m.App.ErrorLog.Println("Cannot get DB repo")
			m.AddErrorMsg(r, "Please login to view categories")
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
		}

		// Get data from form
		tableOrder, _ := strconv.ParseInt(form.Get("table_order"), 10, 64)

		// Update category position
		err = repo.ReorderCategory(int(id), int(tableOrder)+direction)
		if err != nil {
			m.App.ErrorLog.Println(err)
			m.AddErrorMsg(r, "Failed to move category")
			http.Redirect(w, r, "/categories", http.StatusSeeOther)
			return
		}

		// Redirect
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	}
}

func (m *Repository) PostDeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// Get category id from route param
	idParam := chi.URLParam(r, "categoryId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if idParam == "" || err != nil {
		m.AddErrorMsg(r, "Invalid category")
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Get DB repo
	repo, ok := m.GetDB(r)
	if !ok {
		m.App.ErrorLog.Println("Failed to get DB repo")
		m.AddErrorMsg(r, "Log in before deleting categories")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

	// Delete category from database
	err = repo.DeleteCategory(int(id))
	if err != nil {
		m.App.ErrorLog.Println(err)
		if strings.HasPrefix(err.Error(), "cant delete a category that is used") {
			m.AddErrorMsg(r, "Can't delete category that is being used")
		} else {
			m.AddErrorMsg(r, "Failed to delete category")

		}
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Add success message
	m.AddFlashMsg(r, "Category deleted")
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
