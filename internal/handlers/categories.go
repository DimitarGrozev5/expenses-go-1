package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/views/categoriesview"
	"github.com/go-chi/chi"
)

func (m *Repository) Categories(w http.ResponseWriter, r *http.Request) {

	// Get all categories
	categories, err := m.DBClient.GetCategoriesOverview(r.Context(), nil)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting categories")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	// Get time periods
	periods, err := m.DBClient.GetTimePeriods(r.Context(), nil)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting categories")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	// Get user data
	user, err := m.DBClient.GetUser(r.Context(), nil)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Error getting categories")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		return
	}

	// Set default unused categories
	defaultUnused := make([]string, 0, len(categories.Categories))
	for _, category := range categories.Categories {
		defaultUnused = append(defaultUnused, categoryToFormString(category))
	}

	// Set category reset form
	resetForm := forms.NewFromMap(map[string]string{"unused-categories": strings.Join(defaultUnused, ";"), "used-categories": ""})

	// Get template data
	td := models.TemplateData{
		Title: "Categories",
		Form: map[string]*forms.Form{
			"add-category":     forms.New(nil),
			"reset-categories": resetForm,
		},
	}

	// Add forms for expenses
	for _, category := range categories.Categories {
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
		Categories:   categories.Categories,
		TimePeriods:  periods.TimePeriods,
		FreeFunds:    user.FreeFunds,
	}

	// Render view
	data.View().Render(r.Context(), w)
}

func (m *Repository) PostNewCategory(w http.ResponseWriter, r *http.Request) {

	// Parse form
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println(err)
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
	_, err = m.DBClient.AddCategory(r.Context(), &models.AddCategoryParams{
		Name:          name,
		BudgetInput:   budgetInput,
		SpendingLimit: spendingLimit,
		InputInterval: inputInterval,
		InputPeriod:   inputPeriod,
	})
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
			m.App.ErrorLog.Println(err)
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

		// Get data from form
		tableOrder, _ := strconv.ParseInt(form.Get("table_order"), 10, 64)

		// Update category position
		_, err = m.DBClient.ReorderCategory(r.Context(), &models.ReorderCategoryParams{CategoryId: id, NewOrder: tableOrder + int64(direction)})
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
		m.App.ErrorLog.Println(err)
	}

	// Get category id from route param
	idParam := chi.URLParam(r, "categoryId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if idParam == "" || err != nil {
		m.AddErrorMsg(r, "Invalid category")
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Delete category from database
	_, err = m.DBClient.DeleteCategory(r.Context(), &models.DeleteCategoryParams{ID: id})
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

func (m *Repository) PostResetCategories(w http.ResponseWriter, r *http.Request) {

	// Parse form
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println(err)
	}

	// Get form and validate fields
	form := forms.New(r.PostForm)
	form.Required("used-categories")

	// Get used categories string
	usedCategoriesString := form.Get("used-categories")

	if !form.Valid() || len(usedCategoriesString) == 0 {

		// Push form to session
		m.AddForms(r, map[string]*forms.Form{
			"reset-categories": form,
		})

		// Redirect to categories
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Get categories for reset
	cats := strings.Split(usedCategoriesString, ";")

	// Store reset category data
	resetData := []*models.GrpcResetCategoryData{}

	// Loop through categories
	for _, cat := range cats {
		// Get relevant data
		category, err := formStringToCategory(cat)
		if err != nil {
			// Push form to session
			m.AddForms(r, map[string]*forms.Form{
				"reset-categories": form,
			})

			// Add eror message
			m.AddErrorMsg(r, "Error parsing data")

			// Redirect to categories
			http.Redirect(w, r, "/categories", http.StatusSeeOther)
			return
		}

		data := &models.GrpcResetCategoryData{}

		data.Amount = category.InitialAmount
		data.CategoryId = category.ID
		data.BudgetInput = category.BudgetInput
		data.InputInterval = category.InputInterval
		data.InputPeriod = category.InputPeriodId
		data.SpendingLimit = category.SpendingLimit

		resetData = append(resetData, data)
	}

	// Reset all categories
	_, err = m.DBClient.ResetCategories(r.Context(), &models.ResetCategoriesParams{Catgories: resetData})
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.AddErrorMsg(r, "Failed to reset category")
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Add success message
	m.AddFlashMsg(r, "Categories reset")
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func categoryToFormString(c *models.GrpcCategoryOverview) string {
	return fmt.Sprintf(
		"%d,%s,%f,%d,%d,%s,%f,%d,%d,%f,%f",
		c.ID,
		c.Name,
		c.BudgetInput,
		c.InputInterval,
		c.InputPeriodId,
		c.InputPeriodCaption,
		c.SpendingLimit,
		c.PeriodStart.AsTime().Unix(),
		c.PeriodEnd.AsTime().Unix(),
		c.InitialAmount,
		c.CurrentAmount,
	)
}

func formStringToCategory(c string) (*models.GrpcCategoryOverview, error) {
	// Split string to fields
	fields := strings.Split(c, ",")

	// Create empty category
	category := &models.GrpcCategoryOverview{}

	// Check fields length
	if len(fields) < 11 {
		return category, errors.New("wrong format received from frontend")
	}

	// Get category id
	id, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return category, err
	}
	category.ID = id

	// Get amount
	initialAmount, err := strconv.ParseFloat(fields[9], 64)
	if err != nil {
		return category, err
	}
	category.InitialAmount = initialAmount

	// Get budget input
	budgetInput, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return category, err
	}
	category.BudgetInput = budgetInput

	// Get spending limit
	spendingLimit, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return category, err
	}
	category.SpendingLimit = spendingLimit

	// Get input interval
	inputInterval, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return category, err
	}
	category.InputInterval = inputInterval

	// Get input period
	inputPriodId, err := strconv.ParseInt(fields[4], 10, 64)
	if err != nil {
		return category, err
	}
	category.InputPeriodId = inputPriodId

	return category, nil
}
