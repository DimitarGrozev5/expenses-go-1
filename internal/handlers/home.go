package handlers

import (
	"net/http"

	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/views/homeview"

	"github.com/justinas/nosurf"
)

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	// Setup page data
	data := homeview.HomeData{
		Title:     "Home",
		CsrfToken: nosurf.Token(r),
		LoginForm: *forms.New(nil),
	}

	data.View().Render(r.Context(), w)

	// render.RenderTemplate(w, r, "home.page.htm", &models.TemplateData{})
	// views.Test().Render(r.Context(), w)
	// layout.MainLayout("").Render(r.Context(), w)
}
