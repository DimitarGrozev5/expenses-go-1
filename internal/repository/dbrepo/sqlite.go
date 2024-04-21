package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
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

	// Define expensesMap map and slice
	expensesMap := map[int]*models.Expense{}
	expensesOrder := make([]int, 0)

	// Scan rows
	for rows.Next() {
		// Define base models
		var expense models.Expense
		var tag models.Tag

		err = rows.Scan(&expense.ID, &expense.Amount, &expense.Date, &tag.ID, &tag.Name, &tag.UsageCount)
		if err != nil {
			return nil, err
		}

		// Get expense
		oldExpense, ok := expensesMap[expense.ID]

		// If expense hasn't been added
		if !ok {
			expense.Tags = []models.Tag{tag}
			expensesMap[expense.ID] = &expense
			expensesOrder = append(expensesOrder, expense.ID)
			continue
		}

		// If expense has been added
		oldExpense.Tags = append(oldExpense.Tags, tag)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Get expenses slice
	expenses := make([]models.Expense, 0, len(expensesOrder))
	for _, id := range expensesOrder {
		expenses = append(expenses, *expensesMap[id])
	}

	return expenses, nil
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
	exisitingTags := make([]models.Tag, 0, len(expense.Tags))

	// Go through tags
	for _, tag := range expense.Tags {

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
			return err
		}
		defer rows.Close()

		// Scan rows
		for rows.Next() {
			// Define base model
			var tag models.Tag

			err = rows.Scan(&tag.ID, &tag.Name, &tag.UsageCount)
			if err != nil {
				return err
			}

			// Add tag to existing tags
			exisitingTags = append(exisitingTags, tag)
		}
		err = rows.Err()
		if err != nil {
			return err
		}
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

	// Store VALUES template
	tagValuesTmpl := make([]string, 0, len(exisitingTags))

	// Store values
	tagValues := make([]interface{}, 0, len(exisitingTags)*2)

	// Loop trough new tags
	for i, tag := range exisitingTags {

		// Define template
		tmpl := fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)

		// Add to templates
		tagValuesTmpl = append(tagValuesTmpl, tmpl)

		// Add tp values
		tagValues = append(tagValues, int(expenseId), tag.ID)
	}

	// Define query to insert relations
	stmt = `INSERT INTO expense_tags(expense_id, tag_id) VALUES `

	// Append templates
	stmt = fmt.Sprintf("%s%s", stmt, strings.Join(tagValuesTmpl, ","))

	// Insert relations
	_, err = tx.ExecContext(
		ctx,
		stmt,
		tagValues...,
	)
	if err != nil {
		return err
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
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// // Define query
	// stmt := `DELETE FROM expenses WHERE id=$1`

	// // Execute query
	// _, err := m.DB.ExecContext(
	// 	ctx,
	// 	stmt,
	// 	id,
	// )

	// if err != nil {
	// 	return err
	// }

	return nil
}
