package rpcserver

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *DatabaseServer) GetExpenses(ectx context.Context, params *models.GrpcEmpty) (*models.GetExpensesReturns, error) {
	// Get db
	db, ok := m.GetDB(ectx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.GetExpenses(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
