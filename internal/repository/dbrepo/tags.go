package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

// Get all tags ordered by most used
func (m *sqliteDBRepo) GetTags(empty *models.GrpcEmpty) (*models.GetTagsReturns, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT id, name, usage_count FROM tags ORDER BY usage_count DESC;`

	// Execute query
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define tags slice
	tags := make([]*models.GrpcTag, 0)

	// Scan rows
	for rows.Next() {
		tag := &models.GrpcTag{}

		err = rows.Scan(&tag.ID, &tag.Name, &tag.UsageCount)
		if err != nil {
			log.Fatal(err)
		}

		// Add to slice
		tags = append(tags, tag)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &models.GetTagsReturns{Tags: tags}, nil
}

// Update multiple tags
func (m *sqliteDBRepo) UpdateTags(tags []string, etx *sql.Tx) ([]models.Tag, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	var tx *sql.Tx
	if etx != nil {
		tx = etx
	} else {
		var err error
		tx, err = m.DB.Begin()
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()
	}

	// There must be tags
	if len(tags) < 1 {
		return nil, fmt.Errorf("you must have at least one tag")
	}

	// Store VALUES template
	tagValuesTmpl := make([]string, 0, len(tags))

	// Store values
	tagValues := make([]interface{}, 0, len(tags)*3)

	// Loop trough new tags
	for i, tag := range tags {

		// Define template
		tmpl := fmt.Sprintf("($%d)", i+1)

		// Add to templates
		tagValuesTmpl = append(tagValuesTmpl, tmpl)

		// Add tp values
		tagValues = append(tagValues, tag)
	}

	// Define query
	stmt := `INSERT INTO procedure_insert_tag (name) VALUES `

	// Append templates
	stmt = fmt.Sprintf("%s%s", stmt, strings.Join(tagValuesTmpl, ","))

	// Insert tags
	_, err := tx.ExecContext(
		ctx,
		stmt,
		tagValues...,
	)
	if err != nil {
		return nil, err
	}

	// Wrap tags in quotes
	wrapedTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		wrapedTags = append(wrapedTags, fmt.Sprintf("'%s'", tag))
	}

	// Set query
	query := fmt.Sprintf("SELECT id, name, usage_count FROM tags WHERE name IN (%s);", strings.Join(wrapedTags, ", "))

	// Get rows
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allTags []models.Tag

	// Scan rows
	for rows.Next() {
		// Define base model
		var tag models.Tag

		err = rows.Scan(&tag.ID, &tag.Name, &tag.UsageCount)
		if err != nil {
			return nil, err
		}

		// Add tag to existing tags
		allTags = append(allTags, tag)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if etx == nil {
		tx.Commit()
	}

	// Return all tags
	return allTags, nil
}
