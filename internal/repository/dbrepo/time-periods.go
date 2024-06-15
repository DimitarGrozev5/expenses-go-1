package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Get all tags ordered by most used
func (m *sqliteDBRepo) GetTimePeriods(empty *models.GrpcEmpty) (*models.GetTimePeriodsReturns, error) {
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
	periods := make([]*models.GrpcTimePeriod, 0)
	var createdAt time.Time
	var updatedAt sql.NullTime

	// Scan rows
	for rows.Next() {
		period := &models.GrpcTimePeriod{}

		err = rows.Scan(
			&period.ID,
			&period.Period,
			&period.Caption,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			log.Fatal(err)
		}

		period.CreatedAt = timestamppb.New(createdAt)
		if updatedAt.Valid {
			period.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		// Add to slice
		periods = append(periods, period)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &models.GetTimePeriodsReturns{TimePeriods: periods}, nil
}
