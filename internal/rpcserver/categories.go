package rpcserver

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *DatabaseServer) GetCategories(ectx context.Context, params *models.GrpcEmpty) (*models.GetCategoriesReturns, error) {
	// Get db
	db, ok := m.GetDB(ectx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.GetCategories(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
