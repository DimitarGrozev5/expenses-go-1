package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/forms"
	"github.com/dimitargrozev5/expenses-go-1/internal/handlers"
	"github.com/dimitargrozev5/expenses-go-1/internal/helpers"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	seed   = flag.Bool("seed", false, "Create and seed new DB asd@asd.asd with password asd")
	port   = flag.String("port", "3001", "Set server port")
	dbAddr = flag.String("dbaddr", "127.0.0.1:3002", "Database Controller address")
)

// Init app config
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	// Start gRPC client
	var opts = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor),
	}

	conn, err := grpc.NewClient(*dbAddr, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = run(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer handlers.Repo.CloseAllConnections()

	fmt.Println("Starting server on port ", *port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", *port),
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run(conn *grpc.ClientConn) error {
	// res, _ := bcrypt.GenerateFromPassword([]byte("asd"), bcrypt.DefaultCost)
	// fmt.Println(string(res))

	// Register models to Session
	gob.Register(models.User{})
	gob.Register(models.Expense{})
	gob.Register(forms.Form{})
	gob.Register(map[string]*forms.Form{})

	// Read command line arguments
	flag.Parse()

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

	// Pass app config
	helpers.NewHelpers(&app)

	// Create gRPC client
	client := models.NewDatabaseClient(conn)

	// Create handlers repo
	repo := handlers.NewRepo(&app, client)
	handlers.NewHandlers(repo)

	// Seed DB
	if *seed {
		Seed(app.DBPath)
	}

	return nil
}
