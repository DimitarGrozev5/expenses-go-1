package main

import (
	"github.com/dimitargrozev5/expenses-go-1/internal/config"
)

// Init app config
var app config.DBControllerConfig

func main() {

	// Setup app
	setupAppState()

	// Setup grpc service
	setupGrpcService()
}
