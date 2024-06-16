package dbnoderpc

import (
	"context"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
)

// Setup data for Service
type DatabaseServer struct {
	models.UnimplementedDatabaseServer

	App *config.DBNodeConfig
}

// Repository used by the RPC commands
var Server *DatabaseServer

// Creates a new repsoitory
func NewService(a *config.DBNodeConfig) *DatabaseServer {
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
	db, ok := m.App.DBRepos[userKey]

	return db, ok
}

// Get user connection
func (m *DatabaseServer) GetDBConn(ctx context.Context) (*driver.DB, bool) {
	// Get user key
	userKey, ok := ctx.Value("userKey").(string)
	if !ok {
		return nil, false
	}

	// Get db connection
	db, ok := m.App.DBConnections[userKey]

	return db, ok
}
