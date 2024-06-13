package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (m *sqliteDBRepo) GetAccounts(params *models.GetAccountsParams) (*models.GetAccountsReturns, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT id, name, current_amount, usage_count, table_order, created_at, updated_at FROM accounts`

	// Order by option
	if params.OrderByPopularity {
		query = fmt.Sprintf("%s ORDER BY usage_count DESC", query)
	} else {
		query = fmt.Sprintf("%s ORDER BY table_order DESC", query)
	}

	// Get rows
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Set variable for accounts
	accounts := make([]*models.GrpcAccount, 0)

	// Scan rows
	for rows.Next() {
		// Define base models
		account := &models.GrpcAccount{}

		// Get time
		var createdAt time.Time
		var updatedAt time.Time

		err = rows.Scan(
			&account.ID,
			&account.Name,
			&account.CurrentAmount,
			&account.UsageCount,
			&account.TableOrder,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		account.CreatedAt = timestamppb.New(createdAt)
		account.UpdatedAt = timestamppb.New(updatedAt)

		// Add to accounts
		accounts = append(accounts, account)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &models.GetAccountsReturns{Accounts: accounts}, nil
}

func (m *sqliteDBRepo) GetAccount(id int64) (*models.GrpcAccount, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT id, name, current_amount, usage_count, table_order, created_at, updated_at FROM accounts WHERE id=$1`

	// Get row
	row := m.DB.QueryRowContext(ctx, query, id)

	// Set account
	account := &models.GrpcAccount{}
	var createdAt time.Time
	var updatedAt time.Time

	// Scan row into model
	err := row.Scan(
		&account.ID,
		&account.Name,
		&account.CurrentAmount,
		&account.UsageCount,
		&account.TableOrder,
		&createdAt,
		&updatedAt,
	)

	// Check for error
	if err != nil {
		return account, err
	}

	account.CreatedAt = timestamppb.New(createdAt)
	account.UpdatedAt = timestamppb.New(updatedAt)

	return account, nil
}

func (m *sqliteDBRepo) AddAccount(params *models.AddAccountParams) (*models.GrpcEmpty, error) {
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
	stmt := `INSERT INTO procedure_insert_account (name) VALUES ($1)`

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		stmt,
		params.Name,
	)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return nil, nil
}

func (m *sqliteDBRepo) EditAccountName(params *models.EditAccountNameParams) (*models.GrpcEmpty, error) {
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
	stmt := `UPDATE procedure_account_update_name SET name = $1 WHERE id = $2`

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		stmt,
		params.Name,
		params.ID,
	)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return nil, nil
}

func (m *sqliteDBRepo) DeleteAccount(params *models.DeleteAccountParams) (*models.GrpcEmpty, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get account
	account, err := m.GetAccount(params.ID)
	if err != nil {
		return nil, err
	}

	// If account is connected to expenses, don't delete it
	if account.UsageCount > 0 {
		return nil, nil
	}

	// Setup query to delete account
	stmt := `DELETE FROM accounts WHERE id=$1`

	// Execute query
	_, err = tx.ExecContext(ctx, stmt, params.ID)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return nil, nil
}

func (m *sqliteDBRepo) TransferFunds(params *models.TransferFundsParams) (*models.GrpcEmpty, error) {
	return nil, nil
}

func (m *sqliteDBRepo) ReorderAccount(params *models.ReorderAccountParams) (*models.GrpcEmpty, error) {
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
	stmt := `UPDATE procedure_change_accounts_order SET table_order = $1 WHERE id = $2`

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		stmt,
		params.Account.TableOrder+params.Direction,
		params.Account.ID,
	)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return nil, nil
}

// func (m *sqliteDBRepo) GetAccounts(orderByPopularity bool) ([]models.Account, error) {
// 	// Define context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Define query
// 	query := `SELECT id, name, current_amount, usage_count, table_order, created_at, updated_at FROM accounts`

// 	// Order by option
// 	if orderByPopularity {
// 		query = fmt.Sprintf("%s ORDER BY usage_count DESC", query)
// 	} else {
// 		query = fmt.Sprintf("%s ORDER BY table_order DESC", query)
// 	}

// 	// Get rows
// 	rows, err := m.DB.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// Set variable for accounts
// 	accounts := make([]models.Account, 0)

// 	// Scan rows
// 	for rows.Next() {
// 		// Define base models
// 		var account models.Account

// 		err = rows.Scan(
// 			&account.ID,
// 			&account.Name,
// 			&account.CurrentAmount,
// 			&account.UsageCount,
// 			&account.TableOrder,
// 			&account.CreatedAt,
// 			&account.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Add to accounts
// 		accounts = append(accounts, account)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return accounts, nil
// }

// func (m *sqliteDBRepo) GetAccount(id int) (models.Account, error) {
// 	// Define context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Define query
// 	query := `SELECT id, name, current_amount, usage_count, table_order, created_at, updated_at FROM accounts WHERE id=$1`

// 	// Get row
// 	row := m.DB.QueryRowContext(ctx, query, id)

// 	// Set account
// 	var account models.Account

// 	// Scan row into model
// 	err := row.Scan(
// 		&account.ID,
// 		&account.Name,
// 		&account.CurrentAmount,
// 		&account.UsageCount,
// 		&account.TableOrder,
// 		&account.CreatedAt,
// 		&account.UpdatedAt,
// 	)

// 	// Check for error
// 	if err != nil {
// 		return account, err
// 	}

// 	return account, nil
// }

// func (m *sqliteDBRepo) AddAccount(name string) error {
// 	// Define context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Start transaction
// 	tx, err := m.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Define query to insert account
// 	stmt := `INSERT INTO procedure_insert_account (name) VALUES ($1)`

// 	// Execute query
// 	_, err = tx.ExecContext(
// 		ctx,
// 		stmt,
// 		name,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	tx.Commit()
// 	return nil
// }

// func (m *sqliteDBRepo) EditAccountName(id int, name string) error {
// 	// Define context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Start transaction
// 	tx, err := m.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Define query to insert account
// 	stmt := `UPDATE procedure_account_update_name SET name = $1 WHERE id = $2`

// 	// Execute query
// 	_, err = tx.ExecContext(
// 		ctx,
// 		stmt,
// 		name,
// 		id,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	tx.Commit()
// 	return nil
// }

// func (m *sqliteDBRepo) DeleteAccount(id int) error {
// 	// Define context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Start transaction
// 	tx, err := m.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Get account
// 	account, err := m.GetAccount(id)
// 	if err != nil {
// 		return err
// 	}

// 	// Setup query to delete account
// 	stmt := `DELETE FROM accounts WHERE id=$1`

// 	// If account is connected to expenses, don't delete it
// 	if account.UsageCount > 0 {
// 		return nil
// 	}

// 	// Execute query
// 	_, err = tx.ExecContext(ctx, stmt, id)
// 	if err != nil {
// 		return err
// 	}

// 	tx.Commit()
// 	return nil
// }

// func (m *sqliteDBRepo) TransferFunds(fromAccount, toAccount models.Account, amount float64) error {
// 	return nil
// }

// func (m *sqliteDBRepo) ReorderAccount(currentAccount models.Account, direction int) error {
// 	// Define context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Start transaction
// 	tx, err := m.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Setup query
// 	stmt := `UPDATE procedure_change_accounts_order SET table_order = $1 WHERE id = $2`

// 	// Execute query
// 	_, err = tx.ExecContext(
// 		ctx,
// 		stmt,
// 		currentAccount.TableOrder+direction,
// 		currentAccount.ID,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	tx.Commit()

// 	return nil
// }
