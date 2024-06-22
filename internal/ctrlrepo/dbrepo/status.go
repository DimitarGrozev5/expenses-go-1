package dbrepo

import (
	"context"
	"time"
)

func (m sqliteDBRepo) GetVersion() (int64, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get DB Version
	row := m.DB.QueryRowContext(ctx, `PRAGMA user_version`)

	// Pull data from row
	var userVersion int64
	err := row.Scan(&userVersion)
	if err != nil {
		return 0, err
	}

	return userVersion, nil
}
