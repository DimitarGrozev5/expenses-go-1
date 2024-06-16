package main

import "github.com/dimitargrozev5/expenses-go-1/internal/config"

var app config.DBNodeConfig

func main() {

	// Setup app state
	setupAppState()

	// Register with db controller
	registerDBNode()

	// Start gRPC server
	setupGrpcService()
}
