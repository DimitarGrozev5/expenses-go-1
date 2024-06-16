package rpcserver

import (
	"context"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

// Stub that does nothing
func (m *DatabaseServer) RegisterNode(ctx context.Context, params *models.DBNodeData) (*models.GrpcEmpty, error) {
	return nil, nil
}
