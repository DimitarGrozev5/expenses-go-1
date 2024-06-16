package dbnoderpc

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *DatabaseServer) RegisterNode(ctx context.Context, params *models.DBNodeData) (*models.GrpcEmpty, error) {
	// // Get db
	// db, ok := m.GetDB(ctx)
	// if !ok {
	// 	return nil, fmt.Errorf("can't find user db connection")
	// }

	// ret, err := db.GetAccounts(params)
	// if err != nil {
	// 	return nil, err
	// }

	fmt.Println(params)

	return nil, nil
}
