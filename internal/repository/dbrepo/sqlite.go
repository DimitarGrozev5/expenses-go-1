package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Close connection
func (m *sqliteDBRepo) Close() error {
	return m.DB.Close()
}

// Get user by id
func (m *sqliteDBRepo) GetUserByEmail(email string) (models.User, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define empty model
	var u models.User

	// Define query
	query := `SELECT id, email, password
				FROM user WHERE email = $1`

	// Get row
	row := m.DB.QueryRowContext(ctx, query, email)

	// Scan row into model
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
	)

	// Check for error
	if err != nil {
		return u, err
	}

	// Return user
	return u, nil
}

// Authenticate user
func (m *sqliteDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	// Define variables
	var id int
	var hashedPassword string

	// Get user
	u, err := m.GetUserByEmail(email)
	if err != nil {
		return 0, "", err
	}

	// Set variable
	id = u.ID
	hashedPassword = u.Password

	// Check if password matches
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// Get all tags ordered by most used
func (m *sqliteDBRepo) GetTags() ([]models.Tag, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT id, name, usage_count, last_used FROM tags ORDER BY usage_count DESC;`

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

		err = rows.Scan(&tag.ID, &tag.Name, &tag.UsageCount, &tag.LastUsed)
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

// Get expenses ordered by date
func (m *sqliteDBRepo) GetExpenses() ([]models.Expense, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `	SELECT	expenses.id as expense_id,
						expenses.amount,
						expenses.date,
						tags.id as tag_id,
						tags.name as tag_name,
						tags.usage_count
				FROM expenses
				JOIN expense_tags	ON (expenses.id = expense_tags.expense_id)
				JOIN tags			ON (expense_tags.tag_id = tags.id)
				ORDER BY expenses.date DESC, tags.usage_count DESC;`

	// Get rows
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define exenses map
	expenses := map[int]models.Expense{}

	// Scan rows
	for rows.Next() {
		var expense models.Expense
		var tag models.Tag

		err = rows.Scan(&expense.ID, &expense.Amount, &expense.Date, &tag.ID, &tag.Name, &tag.UsageCount)
		if err != nil {
			log.Fatal(err)
		}

		// Get expense
		oldExpense, ok := expenses[expense.ID]
		fmt.Println(oldExpense, ok)
		if !ok {
			expense.Tags = append(expense.Tags, tag)
			expenses[expense.ID] = expense
			continue
		}
		oldExpense.Tags = append(oldExpense.Tags, tag)
		expenses[expense.ID] = oldExpense
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Get expenses slice
	expensesSlice := []models.Expense{}
	for _, expense := range expenses {
		expensesSlice = append(expensesSlice, expense)
	}

	return expensesSlice, nil
}

// Add expense
func (m *sqliteDBRepo) AddExpense(expense models.Expense) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// There must be tags
	if len(expense.Tags) < 1 {
		return fmt.Errorf("you must have at least one tag")
	}

	// Divide tags in to old and new
	newTagsData := make([]models.Tag, 0, len(expense.Tags)*2/3)
	oldTagsData := make([]models.Tag, 0, len(expense.Tags))

	// Go through tags
	for _, tag := range expense.Tags {

		// If tag id is not set, it's a new tag
		if tag.ID == -1 {
			newTagsData = append(newTagsData, tag)
		} else {
			oldTagsData = append(oldTagsData, tag)
		}
	}

	// Insert new tags
	// Define base query
	baseInsertStmt := `INSERT INTO tags(name, usage_count, last_used) VALUES ($1, $2, $3)`

	// Loop trough new tags and insert them
	for _, tag := range newTagsData {

		// Insert tag
		result, err := tx.ExecContext(
			ctx,
			baseInsertStmt,
			tag.Name,
			1,
			time.Now(),
		)
		if err != nil {
			return err
		}

		// Get new tag id
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Append new tag to old tags
		oldTagsData = append(oldTagsData, models.Tag{ID: int(id)})
	}

	// Define query to insert expense
	stmt := `INSERT INTO expenses(amount, date) VALUES($1, $2)`

	// Execute query
	result, err := tx.ExecContext(
		ctx,
		stmt,
		expense.Amount,
		expense.Date,
	)
	if err != nil {
		return err
	}

	// Get new expense expenseId
	expenseId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Define query to insert expense tags
	stmt = `INSERT INTO expense_tags(expense_id, tag_id) VALUES($1, $2)`

	// Loop through old tags
	for _, tag := range oldTagsData {

		// Execute query
		_, err = tx.ExecContext(
			ctx,
			stmt,
			expenseId,
			tag.ID,
		)

		if err != nil {
			return err
		}
	}

	tx.Commit()

	return nil
}

// Edit expense
func (m *sqliteDBRepo) EditExpense(expense models.Expense) error {
	return nil
	// // Define context with timeout
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// // Define query
	// stmt := `UPDATE expenses SET
	// 			amount=$1,
	// 			label=$2,
	// 			date=$3
	// 		WHERE id=$4`

	// // Execute query
	// _, err := m.DB.ExecContext(
	// 	ctx,
	// 	stmt,
	// 	expense.Amount,
	// 	expense.Label,
	// 	expense.Date,
	// 	expense.ID,
	// )

	// if err != nil {
	// 	return err
	// }

	// return nil
}

// Delete expense
func (m *sqliteDBRepo) DeleteExpense(id int) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	stmt := `DELETE FROM expenses WHERE id=$1`

	// Execute query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
