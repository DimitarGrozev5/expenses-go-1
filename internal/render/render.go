package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	// Create template cache
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	// Render template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Get all files
	pages, err := filepath.Glob("./templates/*.page.htm")
	if err != nil {
		return cache, err
	}

	layouts, err := filepath.Glob("./templates/*.layout.htm")
	if err != nil {
		return cache, err
	}

	// Range trough all files
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.htm")
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
