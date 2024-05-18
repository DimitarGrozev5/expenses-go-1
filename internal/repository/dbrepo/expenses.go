package dbrepo

import (
	"context"
	"database/sql"
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
						tags.usage_count,
						accounts.id as account_id,
						accounts.name as account_name
				FROM expenses
				JOIN expense_tags	ON (expenses.id = expense_tags.expense_id)
				JOIN tags			ON (expense_tags.tag_id = tags.id)
				JOIN accounts		ON (expenses.from_account = accounts.id)
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
		var account models.Account

		err = rows.Scan(&expense.ID, &expense.Amount, &expense.Date, &tag.ID, &tag.Name, &tag.UsageCount, &account.ID, &account.Name)
		if err != nil {
			return nil, err
		}

		// Get expense
		oldExpense, ok := expensesMap[expense.ID]

		// If expense hasn't been added
		if !ok {
			expense.Tags = []models.Tag{tag}
			expense.FromAccount = account
			expensesMap[expense.ID] = &expense
			expensesOrder = append(expensesOrder, expense.ID)
			continue
		}

		// If expense has been added
		oldExpense.Tags = append(oldExpense.Tags, tag)
		oldExpense.FromAccount = account
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
func (m *sqliteDBRepo) AddExpense(expense models.Expense, tags []string) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update tags
	exisitingTags, err := m.UpdateTags(tags, nil)
	if err != nil {
		return err
	}

	// Define query to insert expense
	stmt := `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES($1, $2, $3, $4)`

	// Execute query
	result, err := tx.ExecContext(
		ctx,
		stmt,
		expense.Amount,
		expense.Date,
		expense.FromAccountId,
		expense.FromCategoryId,
	)
	if err != nil {
		return err
	}

	// Get new expense expenseId
	expenseId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(expenseId)
	fmt.Println(exisitingTags)

	// Add tag relations
	// err = m.AddExpenseTags(int(expenseId), exisitingTags, tx)
	// if err != nil {
	// 	return err
	// }

	tx.Commit()

	return nil
}

// Edit expense
func (m *sqliteDBRepo) EditExpense(expense models.Expense, tags []string) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create query to remove all expense tags
	stmt := `DELETE FROM procedure_unlink_tags_from_expense WHERE expense_id = $1;`

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
	allTags, err := m.UpdateTags(tags, tx)
	if err != nil {
		return err
	}

	// Update expense
	stmt = `UPDATE procedure_update_expense SET
				amount = $1,
				date = $2,
				from_account = $3,
				from_category = $4
			WHERE id = $5`

	// Update expense
	_, err = tx.ExecContext(
		ctx,
		stmt,
		expense.Amount,
		expense.Date,
		expense.FromAccountId,
		expense.FromCategoryId,
		expense.ID,
	)
	if err != nil {
		return err
	}

	// Add relations
	err = m.AddExpenseTags(expense.ID, allTags, tx)
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
	stmt := `DELETE FROM procedure_remove_expense WHERE id=$1`

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
		defer tx.Rollback()
	}

	// Store VALUES template
	tagValuesTmpl := make([]string, 0, len(tags))

	// Store values
	tagValues := make([]interface{}, 0, len(tags)*2)

	// Loop trough new tags
	for i, tag := range tags {

		// Define template
		tmpl := fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)

		// Add to templates
		tagValuesTmpl = append(tagValuesTmpl, tmpl)

		// Add tp values
		tagValues = append(tagValues, expenseId, tag.ID)
	}

	// Define query to insert relations
	stmt := `INSERT INTO procedure_link_tag_to_expense(expense_id, tag_id) VALUES `

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

	if etx == nil {
		tx.Commit()
	}

	return nil
}
