package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/handlers"
	"github.com/dimitargrozev5/expenses-go-1/internal/helpers"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/render"
)

const portNumber = ":8080"

// Init app config
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer handlers.Repo.CloseAllConnections()

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

func run() error {
	// res, _ := bcrypt.GenerateFromPassword([]byte("asdasd123"), bcrypt.DefaultCost)
	// fmt.Println(string(res))

	// Register models to Session
	gob.Register(models.User{})
	gob.Register(models.Expense{})

	// Set in production
	app.InProduction = false

	// Set info log
	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// Set error log
	errorLog = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Set session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Set db path
	app.DBPath = "./db/"

	// Create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false

	// Pass app config
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	// Create handlers repo
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	return nil
}
