package dbrepo

import (
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *sqliteDBRepo) RegisterNode(params *models.DBNodeData) (*models.GrpcEmpty, error) {
	// // Define context with timeout
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// // Start transaction
	// tx, err := m.DB.Begin()
	// if err != nil {
	// 	return nil, err
	// }
	// defer tx.Rollback()

	return nil, nil
}
