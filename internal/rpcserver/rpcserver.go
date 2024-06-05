package rpcserver

import (
	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

// Setup data for Service
type DatabaseServer struct {
	models.UnimplementedDatabaseServer

	App *config.DBControllerConfig
}

// Repository used by the RPC commands
var Server *DatabaseServer

// Creates a new repsoitory
func NewService(a *config.DBControllerConfig) *DatabaseServer {
	return &DatabaseServer{
		App: a,
	}
}

// Sets the repository for the handlers
func NewDatabaseServer(r *DatabaseServer) {
	Server = r
}
