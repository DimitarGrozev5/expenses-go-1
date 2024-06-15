package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (m *sqliteDBRepo) GetCategories(params *models.GrpcEmpty) (*models.GetCategoriesReturns, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT
				id,
				name,
				budget_input,
				last_input_date,
				next_input_date,
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
	categories := make([]*models.GrpcCategory, 0)

	// Scan rows
	for rows.Next() {
		// Define base models
		category := &models.GrpcCategory{}
		var lastInputDate time.Time
		var createdAt time.Time
		var updatedAt sql.NullTime

		// Store duration
		var nextInputDate string

		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.BudgetInput,
			&lastInputDate,
			&nextInputDate,
			&category.SpendingLimit,
			&category.SpendingLeft,
			&category.InitialAmount,
			&category.CurrentAmount,
			&category.TableOrder,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse duration
		// t, err := time.Parse("2006-01-02 15:04:05", nextInputDate)
		// if err != nil {
		// 	return nil, err
		// }
		// category.InputInterval = t.Sub(category.LastInputDate)

		category.LastInputDate = timestamppb.New(lastInputDate)
		category.CreatedAt = timestamppb.New(createdAt)
		if updatedAt.Valid {
			category.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		// Add to accounts
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &models.GetCategoriesReturns{Categories: categories}, nil
}

func (m *sqliteDBRepo) GetCategoriesOverview(params *models.GrpcEmpty) (*models.GetCategoriesOverviewReturns, error) {
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
	categories := make([]*models.GrpcCategoryOverview, 0)

	// Scan rows
	for rows.Next() {
		// Define base models
		category := &models.GrpcCategoryOverview{}
		var periodStart time.Time
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
			&periodStart,
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
		category.PeriodStart = timestamppb.New(periodStart)
		category.PeriodEnd = timestamppb.New(t)

		// Add to accounts
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &models.GetCategoriesOverviewReturns{Categories: categories}, nil
}

func (m *sqliteDBRepo) AddCategory(params *models.AddCategoryParams) (*models.GrpcEmpty, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

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
	_, err = tx.ExecContext(
		ctx,
		stmt,
		params.Name,
		params.BudgetInput,
		params.InputInterval,
		params.InputPeriod,
		params.SpendingLimit,
	)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return nil, nil
}

func (m *sqliteDBRepo) DeleteCategory(params *models.DeleteCategoryParams) (*models.GrpcEmpty, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Setup query to delete account
	stmt := `DELETE FROM procedure_remove_category WHERE id=$1`

	// Execute query
	_, err = tx.ExecContext(ctx, stmt, params.ID)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return nil, nil
}

func (m *sqliteDBRepo) ReorderCategory(params *models.ReorderCategoryParams) (*models.GrpcEmpty, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Setup query
	stmt := `UPDATE procedure_change_categories_order SET table_order = $1 WHERE id = $2`

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		stmt,
		params.NewOrder,
		params.CategoryId,
	)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return nil, nil
}

func (m *sqliteDBRepo) ResetCategory(amount float64, categoryId int64, budgetInput float64, inputInterval int64, inputPeriod int64, spendingLimit float64, etx *sql.Tx) error {
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
	_, err := tx.ExecContext(
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

	if etx == nil {
		tx.Commit()
	}

	return nil
}

func (m *sqliteDBRepo) ResetCategories(params *models.ResetCategoriesParams) (*models.GrpcEmpty, error) {

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, categoryData := range params.Catgories {
		err = m.ResetCategory(categoryData.Amount, categoryData.CategoryId, categoryData.BudgetInput, categoryData.InputInterval, categoryData.InputPeriod, categoryData.SpendingLimit, tx)
		if err != nil {
			return nil, err
		}
	}

	tx.Commit()
	return nil, nil
}
