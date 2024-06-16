package dbnoderpc

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *DatabaseServer) GetAccounts(ctx context.Context, params *models.GetAccountsParams) (*models.GetAccountsReturns, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.GetAccounts(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) AddAccount(ctx context.Context, params *models.AddAccountParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.AddAccount(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) EditAccountName(ctx context.Context, params *models.EditAccountNameParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.EditAccountName(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) DeleteAccount(ctx context.Context, params *models.DeleteAccountParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.DeleteAccount(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) TransferFunds(ctx context.Context, params *models.TransferFundsParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.TransferFunds(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) ReorderAccount(ctx context.Context, params *models.ReorderAccountParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.ReorderAccount(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
