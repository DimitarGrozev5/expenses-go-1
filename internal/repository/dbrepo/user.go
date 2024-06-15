package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Close connection
func (m *sqliteDBRepo) Close() error {
	return m.DB.Close()
}

// Get user by id
func (m *sqliteDBRepo) GetUser(empty *models.GrpcEmpty) (*models.GrpcUser, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define empty model
	u := &models.GrpcUser{}
	var createdAt time.Time
	var updatedAt sql.NullTime

	// Define query
	query := `SELECT id, email, password, db_version, free_funds, created_at, updated_at
				FROM user LIMIT 1`

	// Get row
	row := m.DB.QueryRowContext(ctx, query)

	// Scan row into model
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.DBVersion,
		&u.FreeFunds,
		&createdAt,
		&updatedAt,
	)

	// Check for error
	if err != nil {
		return u, err
	}

	u.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		u.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	// Return user
	return u, nil
}

// Authenticate user
func (m *sqliteDBRepo) Authenticate(testPassword string) (int64, string, int64, error) {
	// Get user
	u, err := m.GetUser(nil)
	if err != nil {
		return 0, "", 0, err
	}

	// Set variable
	id := u.ID
	hashedPassword := u.Password
	dbVersion := u.DBVersion

	// Check if password matches
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, "", 0, err
	}

	return id, hashedPassword, dbVersion, nil
}

// Modify free funds
func (m *sqliteDBRepo) ModifyFreeFunds(params *models.ModifyFreeFundsParams) (*models.GrpcEmpty, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update tags
	tags, err := m.UpdateTags([]string{params.TagName}, tx)
	if err != nil || len(tags) != 1 {
		return nil, err
	}

	// Define query to insert account
	stmt := `INSERT INTO procedure_add_free_funds (
		amount,
		to_account,
		tag_id
	) VALUES (
		$1,
		$2,
		$3
	)`

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		stmt,
		params.Amount,
		params.ToAccountId,
		tags[0].ID,
	)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return nil, nil
}
