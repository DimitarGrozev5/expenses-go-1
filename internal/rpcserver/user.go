package rpcserver

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *DatabaseServer) GetUser(ctx context.Context, params *models.GrpcEmpty) (*models.GrpcUser, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.GetUser(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) ModifyFreeFunds(ctx context.Context, params *models.ModifyFreeFundsParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.ModifyFreeFunds(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
