package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dimitargrozev5/expenses-go-1/pkg/config"
	"github.com/dimitargrozev5/expenses-go-1/pkg/handlers"
	"github.com/dimitargrozev5/expenses-go-1/pkg/render"
)

const portNumber = ":8080"

// Init app config
var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// Create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	// Pass app config
	render.NewTemplates(&app)

	// Create handlers repo
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	fmt.Println("Starting server on port ", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
