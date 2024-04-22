package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

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

	// Update tags
	exisitingTags, err := m.UpdateTags(expense.Tags, tx)
	if err != nil {
		return err
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

	// Add tag relations
	err = m.AddExpenseTags(int(expenseId), exisitingTags, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

// Edit expense
func (m *sqliteDBRepo) EditExpense(expense models.Expense) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// Create query to remove all expense tags
	stmt := `DELETE FROM expense_tags WHERE expense_id = $1`

	// Delete relations
	_, err = tx.ExecContext(
		ctx,
		stmt,
		expense.ID,
	)
	if err != nil {
		return err
	}

	// Update tags
	tags, err := m.UpdateTags(expense.Tags, tx)
	if err != nil {
		return err
	}

	// Update expense
	stmt = `UPDATE expenses SET
				amount = $1,
				date = $2,
				updated_at = $3
			WHERE id = $4`

	// Update expense
	_, err = tx.ExecContext(
		ctx,
		stmt,
		expense.Amount,
		expense.Date,
		time.Now(),
		expense.ID,
	)
	if err != nil {
		return err
	}

	// Add relations
	err = m.AddExpenseTags(expense.ID, tags, tx)
	if err != nil {
		return err
	}

	// Commit to transaction and exit
	tx.Commit()
	return nil
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

// Add relations based on tags
func (m *sqliteDBRepo) AddExpenseTags(expenseId int, tags []models.Tag, etx *sql.Tx) error {
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
			return err
		}
	}

	// Store VALUES template
	tagValuesTmpl := make([]string, 0, len(tags))

	// Store values
	tagValues := make([]interface{}, 0, len(tags)*2)

	// Loop trough new tags
	for i, tag := range tags {

		// If tag is new
		if tag.ID == -1 {
			return errors.New("new tag found in tags list, can't add new tag to expense")
		}

		// Define template
		tmpl := fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)

		// Add to templates
		tagValuesTmpl = append(tagValuesTmpl, tmpl)

		// Add tp values
		tagValues = append(tagValues, expenseId, tag.ID)
	}

	// Define query to insert relations
	stmt := `INSERT INTO expense_tags(expense_id, tag_id) VALUES `

	// Append templates
	stmt = fmt.Sprintf("%s%s", stmt, strings.Join(tagValuesTmpl, ","))

	// Insert relations
	_, err := tx.ExecContext(
		ctx,
		stmt,
		tagValues...,
	)
	if err != nil {
		return err
	}

	return nil
}
