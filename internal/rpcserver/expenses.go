package rpcserver

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *DatabaseServer) GetExpenses(ctx context.Context, params *models.GrpcEmpty) (*models.GetExpensesReturns, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.GetExpenses(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) AddExpense(ctx context.Context, params *models.ExpensesParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.AddExpense(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) EditExpense(ctx context.Context, params *models.ExpensesParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.EditExpense(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) DeleteExpense(ctx context.Context, params *models.DeleteExpenseParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.DeleteExpense(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
