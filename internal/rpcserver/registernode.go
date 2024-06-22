package rpcserver

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

// Stub that does nothing
func (m *DatabaseServer) RegisterNode(ctx context.Context, params *models.DBNodeData) (*models.GrpcEmpty, error) {
	fmt.Println(params)

	// Get db
	db := m.App.CtrlDBRepo

	_, err := db.RegisterNode(params)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
