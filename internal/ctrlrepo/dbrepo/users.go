package dbrepo

import (
	"context"
	"database/sql"
	"time"
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

func (m sqliteDBRepo) AddNewUser(email string, version int64) error {
	return nil
}
