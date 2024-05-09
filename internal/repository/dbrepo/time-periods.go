package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

// Get all tags ordered by most used
func (m *sqliteDBRepo) GetTimePeriods() ([]models.TimePeriod, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT id, period, caption, created_at, updated_at FROM time_periods;`

	// Execute query
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define periods slice
	periods := make([]models.TimePeriod, 0)

	// Scan rows
	for rows.Next() {
		var period models.TimePeriod

		err = rows.Scan(
			&period.ID,
			&period.Period,
			&period.Caption,
			&period.CreatedAt,
			&period.UpdatedAt,
		)
		if err != nil {
			log.Fatal(err)
		}

		// Add to slice
		periods = append(periods, period)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return periods, nil
}
