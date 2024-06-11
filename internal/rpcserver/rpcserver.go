package rpcserver

import (
	"context"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
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

// Get user connection
func (m *DatabaseServer) GetDB(ctx context.Context) (repository.DatabaseRepo, bool) {
	// Get user key
	userKey, ok := ctx.Value("userKey").(string)
	if !ok {
		return nil, false
	}

	// Get db connection
	db, ok := m.App.DBConnections[userKey]

	return db, ok
}
