package rpcserver

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *DatabaseServer) GetAccounts(ectx context.Context, params *models.GetAccountsParams) (*models.GetAccountsReturns, error) {
	// Get db
	db, ok := m.GetDB(ectx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.GetAccounts(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
