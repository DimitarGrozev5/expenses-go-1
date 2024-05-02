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
func (m *sqliteDBRepo) GetTags() ([]models.Tag, error) {
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
	tags := make([]models.Tag, 0)

	// Scan rows
	for rows.Next() {
		var tag models.Tag

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

	return tags, nil
}

// Update multiple tags
func (m *sqliteDBRepo) UpdateTags(tags []models.Tag, etx *sql.Tx) ([]models.Tag, error) {
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

	// Divide tags in to old and new
	newTagsData := make([]models.Tag, 0, len(tags)*2/3)
	exisitingTags := make([]models.Tag, 0, len(tags))

	// Go through tags
	for _, tag := range tags {

		// If tag id is not set, it's a new tag
		if tag.ID == -1 {
			newTagsData = append(newTagsData, tag)
		} else {
			exisitingTags = append(exisitingTags, tag)
		}
	}

	// If there are new tags, add them to DB
	if len(newTagsData) > 0 {
		// Store VALUES template
		tagValuesTmpl := make([]string, 0, len(newTagsData))

		// Store values
		tagValues := make([]interface{}, 0, len(newTagsData)*3)

		// Loop trough new tags
		for i, tag := range newTagsData {

			// Define template
			tmpl := fmt.Sprintf("($%d)", i+1)

			// Add to templates
			tagValuesTmpl = append(tagValuesTmpl, tmpl)

			// Add tp values
			tagValues = append(tagValues, tag.Name)
		}

		// Define query
		stmt := `INSERT INTO tags(name) VALUES `

		// Append templates
		stmt = fmt.Sprintf("%s%s RETURNING id, name, usage_count", stmt, strings.Join(tagValuesTmpl, ","))

		// Insert tags
		rows, err := tx.QueryContext(
			ctx,
			stmt,
			tagValues...,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// Scan rows
		for rows.Next() {
			// Define base model
			var tag models.Tag

			err = rows.Scan(&tag.ID, &tag.Name, &tag.UsageCount)
			if err != nil {
				return nil, err
			}

			// Add tag to existing tags
			exisitingTags = append(exisitingTags, tag)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}

	// Return all tags
	return exisitingTags, nil
}
