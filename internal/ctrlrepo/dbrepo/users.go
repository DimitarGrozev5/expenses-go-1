package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m sqliteDBRepo) GetMinUserVersion() (int64, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get DB Version
	row := m.DB.QueryRowContext(ctx, `SELECT MIN(db_version) FROM users;`)

	// Pull data from row
	var userVersion sql.NullInt64
	err := row.Scan(&userVersion)
	if err != nil {
		return 0, err
	}

	if userVersion.Valid {
		return userVersion.Int64, nil
	}

	return 0, nil
}

func (m sqliteDBRepo) GetMaxUserVersion() (int64, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get DB Version
	row := m.DB.QueryRowContext(ctx, `SELECT MAX(db_version) FROM users;`)

	// Pull data from row
	var userVersion sql.NullInt64
	err := row.Scan(&userVersion)
	if err != nil {
		return 0, err
	}

	if userVersion.Valid {
		return userVersion.Int64, nil
	}

	return 0, nil
}

func (m sqliteDBRepo) AddNewUser(email string, password string, version int64) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Generate pass hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Set query
	stmt := `INSERT INTO users (user_email, password_hash, db_version, status) VALUES ($1, $2, $3, 1)`

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		stmt,
		email,
		hash,
		version,
	)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
