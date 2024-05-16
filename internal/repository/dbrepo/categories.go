package dbrepo

import (
	"context"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *sqliteDBRepo) GetCategoriesCount() (int, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Set variable for count
	var count int

	// Define query
	query := `SELECT COUNT(*) FROM categories`

	// Get row
	rows := m.DB.QueryRowContext(ctx, query)

	err := rows.Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (m *sqliteDBRepo) GetCategories() ([]models.Category, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT
				id,
				name,
				budget_input,
				last_input_date,
				input_interval,
				spending_limit,
				spending_left,
				initial_amount,
				current_amount,
				table_order,
				created_at,
				updated_at
			FROM view_categories
			ORDER BY table_order DESC;`

	// Get rows
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Set variable for categories
	categories := make([]models.Category, 0)

	// Scan rows
	for rows.Next() {
		// Define base models
		var category models.Category

		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.BudgetInput,
			&category.LastInputDate,
			&category.InputInterval,
			&category.SpendingLimit,
			&category.SpendingLeft,
			&category.InitialAmount,
			&category.CurrentAmount,
			&category.TableOrder,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Add to accounts
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (m *sqliteDBRepo) GetCategoriesOverview() ([]models.CategoryOverview, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT
				id,
				name,
				budget_input,
				input_interval,
				input_period,
				period_caption,
				spending_limit,
				spending_left,
				period_start,
				period_end,
				initial_amount,
				current_amount,
				can_be_deleted,
				table_order
			FROM view_categories_overview
			ORDER BY table_order DESC;`

	// Get rows
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Set variable for categories
	categories := make([]models.CategoryOverview, 0)

	// Scan rows
	for rows.Next() {
		// Define base models
		var category models.CategoryOverview
		var periodEnd string

		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.BudgetInput,
			&category.InputInterval,
			&category.InputPeriodId,
			&category.InputPeriodCaption,
			&category.SpendingLimit,
			&category.SpendingLeft,
			&category.PeriodStart,
			&periodEnd,
			&category.InitialAmount,
			&category.CurrentAmount,
			&category.CanBeDeleted,
			&category.TableOrder,
		)
		if err != nil {
			return nil, err
		}

		t, err := time.Parse("2006-01-02 15:04:05", periodEnd)
		if err != nil {
			return nil, err
		}
		category.PeriodEnd = t

		// Add to accounts
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (m *sqliteDBRepo) ResetCategories(amount float64, categoryId int, budgetInput float64, inputInterval int, inputPeriod int, spendingLimit float64) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query to insert account
	stmt := `INSERT INTO procedure_fund_category_and_reset_period (
		amount,
		category,
		budget_input,
		input_interval,
		input_period,
		spending_limit
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6
	)`

	// Execute query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		amount,
		categoryId,
		budgetInput,
		inputInterval,
		inputPeriod,
		spendingLimit,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *sqliteDBRepo) AddCategory(name string, budgetInput float64, spendingLimit float64, inputInterval int, inputPeriod int) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query to insert account
	stmt := `INSERT INTO procedure_new_category (
		name,
		budget_input,
		input_interval,
		input_period,
		spending_limit
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)`

	// Execute query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		name,
		budgetInput,
		inputInterval,
		inputPeriod,
		spendingLimit,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *sqliteDBRepo) EditAccountName1(id int, name string) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query to insert account
	stmt := `UPDATE procedure_account_update_name SET name = $1 WHERE id = $2`

	// Execute query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		name,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *sqliteDBRepo) DeleteCategory(id int) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Setup query to delete account
	stmt := `DELETE FROM procedure_remove_category WHERE id=$1`

	// Execute query
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (m *sqliteDBRepo) TransferFunds1(fromAccount, toAccount models.Account, amount float64) error {
	return nil
}

func (m *sqliteDBRepo) ReorderCategory(categoryid int, new_order int) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Setup query
	stmt := `UPDATE procedure_change_categories_order SET table_order = $1 WHERE id = $2`

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		stmt,
		new_order,
		categoryid,
	)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
