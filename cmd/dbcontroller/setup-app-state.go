package main

import (
	"flag"
	"log"
	"os"

	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
)

var dbConn map[string]repository.DatabaseRepo
var infoLog *log.Logger
var errorLog *log.Logger
var dbPath = flag.String("db-path", "./db/", "Path to folder containing sqlite databases")

// Setup app wide state
func setupAppState() {

	// Parse command line flags
	flag.Parse()

	// Set in production
	app.InProduction = false

	// Set info log
	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// Set error log
	errorLog = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Set DB connections
	app.DBConnections = dbConn

	// Set db path
	app.DBPath = *dbPath
}
