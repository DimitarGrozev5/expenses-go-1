package main

import (
	"github.com/dimitargrozev5/expenses-go-1/internal/config"
)

// Init app config
var app config.DBControllerConfig

func main() {

	// Setup app
	setupAppState()

	// Setup and connect to DB
	setupDb()

	// Close db connection on exit
	defer app.CtrlDB.Close()

	// Setup grpc service
	// setupGrpcService()
}
