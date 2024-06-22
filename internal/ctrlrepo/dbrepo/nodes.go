package dbrepo

import (
	"context"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *sqliteDBRepo) GetNodes() ([]models.DBNode, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT id, remote_address, created_at, updated_at FROM db_nodes;`

	// Get rows
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Set variable for dbNodes
	dbNodes := make([]models.DBNode, 0)

	// Scan rows
	for rows.Next() {
		// Define base models
		dbNode := models.DBNode{}

		err = rows.Scan(
			&dbNode.ID,
			&dbNode.RemoteAddress,
			&dbNode.CreatedAt,
			&dbNode.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Add to nodes
		dbNodes = append(dbNodes, dbNode)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return dbNodes, nil
}

func (m *sqliteDBRepo) NewNode() (int64, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Define query to insert account
	stmt := `INSERT INTO db_nodes (remote_address) VALUES ("");`

	// Execute query
	res, err := tx.ExecContext(ctx, stmt)
	if err != nil {
		return 0, err
	}

	// Get last inserted id
	lid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	tx.Commit()
	return lid, nil
}

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
