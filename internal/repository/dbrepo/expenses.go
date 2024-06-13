package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Get expenses ordered by date
func (m *sqliteDBRepo) GetExpenses(param *models.GrpcEmpty) (*models.GetExpensesReturns, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `	SELECT
					expense_id,
					amount,
					date,

					tag_id,
					tag_name,
					usage_count,

					account_id,
					account_name,

					category_id,
					category_name
				FROM view_detailed_expenses;`

	// Get rows
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define expensesMap map and slice
	expensesMap := map[int64]*models.GrpcExpense{}
	expensesOrder := make([]int64, 0)

	// Scan rows
	for rows.Next() {
		// Define base models
		expense := models.GrpcExpense{}
		tag := models.GrpcTag{}
		account := models.GrpcAccount{}
		category := models.GrpcCategory{}
		var date time.Time

		err = rows.Scan(
			&expense.ID,
			&expense.Amount,
			&date,
			&tag.ID,
			&tag.Name,
			&tag.UsageCount,
			&account.ID,
			&account.Name,
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}

		expense.Date = timestamppb.New(date)

		// Get expense
		oldExpense, ok := expensesMap[expense.ID]

		// If expense hasn't been added
		if !ok {
			expense.Tags = []*models.GrpcTag{&tag}
			expense.FromAccount = &account
			expense.FromCategory = &category
			expensesMap[expense.ID] = &expense
			expensesOrder = append(expensesOrder, expense.ID)
			continue
		}

		// If expense has been added
		oldExpense.Tags = append(oldExpense.Tags, &tag)
		oldExpense.FromAccount = &account
		oldExpense.FromCategory = &category
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Get expenses slice
	expenses := make([]*models.GrpcExpense, 0, len(expensesOrder))
	for _, id := range expensesOrder {
		expenses = append(expenses, expensesMap[id])
	}

	return &models.GetExpensesReturns{Expenses: expenses}, nil
}

// Add expense
func (m *sqliteDBRepo) AddExpense(param *models.ExpensesParams) (*models.GrpcEmpty, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update tags
	exisitingTags, err := m.UpdateTags(param.Tags, tx)
	if err != nil {
		return nil, err
	}

	// Define query to insert expense
	// Transactions lock the db so the last expense will be the one inserted
	// Autoincrement adds one so the larges id will be the last
	// Can't use RETURNING because expense insert happens through a trigger
	stmt := `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES($1, $2, $3, $4);`

	// Exec statement
	_, err = tx.ExecContext(
		ctx,
		stmt,
		param.Expense.Amount,
		param.Expense.Date.AsTime(),
		param.Expense.FromAccountId,
		param.Expense.FromCategoryId,
	)
	if err != nil {
		return nil, err
	}

	// Take last inserted expense
	query := `SELECT id FROM expenses ORDER BY id DESC LIMIT 1;`

	// Execute query
	row := tx.QueryRowContext(ctx, query)

	// Get new expense expenseId
	var expenseId int64

	err = row.Scan(&expenseId)

	// Check for error
	if err != nil {
		return nil, err
	}

	// Add tag relations
	err = m.AddExpenseTags(expenseId, exisitingTags, tx)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return nil, nil
}

// Edit expense
func (m *sqliteDBRepo) EditExpense(param *models.ExpensesParams) (*models.GrpcEmpty, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create query to remove all expense tags
	stmt := `DELETE FROM procedure_unlink_tags_from_expense WHERE expense_id = $1;`

	// Delete relations
	_, err = tx.ExecContext(
		ctx,
		stmt,
		param.Expense.ID,
	)
	if err != nil {
		return nil, err
	}

	// Update tags
	allTags, err := m.UpdateTags(param.Tags, tx)
	if err != nil {
		return nil, err
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
		param.Expense.Amount,
		param.Expense.Date.AsTime(),
		param.Expense.FromAccountId,
		param.Expense.FromCategoryId,
		param.Expense.ID,
	)
	if err != nil {
		return nil, err
	}

	// Add relations
	err = m.AddExpenseTags(param.Expense.ID, allTags, tx)
	if err != nil {
		return nil, err
	}

	// Commit to transaction and exit
	tx.Commit()
	return nil, nil
}

// Delete expense
func (m *sqliteDBRepo) DeleteExpense(param *models.DeleteExpenseParams) (*models.GrpcEmpty, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Define query
	stmt := `DELETE FROM procedure_remove_expense WHERE id=$1`

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		stmt,
		param.ID,
	)

	if err != nil {
		return nil, err
	}

	tx.Commit()
	return nil, nil
}

// Add relations based on tags
func (m *sqliteDBRepo) AddExpenseTags(expenseId int64, tags []models.Tag, etx *sql.Tx) error {
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
