package rpcserver

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *DatabaseServer) GetTimePeriods(ctx context.Context, empty *models.GrpcEmpty) (*models.GetTimePeriodsReturns, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.GetTimePeriods(empty)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
