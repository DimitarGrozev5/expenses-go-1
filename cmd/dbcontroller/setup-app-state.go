package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"os"

	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
)

var dbConn = map[string]repository.DatabaseRepo{}
var infoLog *log.Logger
var errorLog *log.Logger
var dbPath = flag.String("db-path", "./db/", "Path to folder containing sqlite databases")
var jwtKeyPath = flag.String("jwt-key-path", "./keys/jwt.pem", "Secret key for signing Json Web Tokens")

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

	// Read the PEM file
	pemData, err := os.ReadFile(*jwtKeyPath)
	if err != nil {
		log.Fatalf("Error reading PEM file: %v", err)
	}

	// Decode the PEM block
	pemBlock, _ := pem.Decode(pemData)
	if pemBlock == nil || pemBlock.Type != "EC PRIVATE KEY" {
		log.Fatalf("Failed to decode PEM block containing private key")
	}

	// Parse the ECDSA private key
	privateKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	if err != nil {
		log.Fatalf("Error parsing ECDSA private key: %v", err)
	}

	// Set jwt key
	app.JWTKey = privateKey
}
