package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *sqliteDBRepo) GetAccountsCount() (int, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT COUNT(*) FROM accounts`

	// Store count
	var count int

	// Get rows
	row := m.DB.QueryRowContext(ctx, query)

	// Scan row into model
	err := row.Scan(&count)

	// Check for error
	if err != nil {
		return count, err
	}

	return count, nil
}

func (m *sqliteDBRepo) GetAccounts(orderByPopularity bool) ([]models.Account, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT id, name, current_amount, usage_count, table_order FROM accounts`

	// Order by option
	if orderByPopularity {
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
	accounts := make([]models.Account, 0)

	// Scan rows
	for rows.Next() {
		// Define base models
		var account models.Account

		err = rows.Scan(&account.ID, &account.Name, &account.CurrentAmount, &account.UsageCount, &account.TableOrder)
		if err != nil {
			return nil, err
		}

		// Add to accounts
		accounts = append(accounts, account)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (m *sqliteDBRepo) GetAccount(id int) (models.Account, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT id, name, current_amount, usage_count, table_order FROM accounts WHERE id=$1`

	// Get row
	row := m.DB.QueryRowContext(ctx, query, id)

	// Set account
	var account models.Account

	// Scan row into model
	err := row.Scan(
		&account.ID,
		&account.Name,
		&account.CurrentAmount,
		&account.UsageCount,
		&account.TableOrder,
	)

	// Check for error
	if err != nil {
		return account, err
	}

	return account, nil
}

func (m *sqliteDBRepo) AddAccount(account models.Account) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query to insert account
	stmt := `INSERT INTO accounts(name) VALUES($1)`

	// Execute query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		account.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *sqliteDBRepo) EditAccount(account models.Account) error {
	return nil
}

func (m *sqliteDBRepo) DeleteAccount(id int) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get account
	account, err := m.GetAccount(id)
	if err != nil {
		return err
	}

	// Setup query to delete account
	stmt := `DELETE FROM accounts WHERE id=$1`

	// If account is connected to expenses, don't delete it
	if account.UsageCount > 0 {
		return nil
	}

	// Execute query
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (m *sqliteDBRepo) TransferFunds(fromAccount, toAccount models.Account, amount float64) error {
	return nil
}

func (m *sqliteDBRepo) ReorderAccount(currentAccount models.Account, direction int) error {
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
	stmt := `
			UPDATE accounts SET table_order = $1 WHERE table_order = $2;
			UPDATE accounts SET table_order = $3 WHERE id = $4;
	`

	// Execute query
	_, err = tx.ExecContext(
		ctx,
		stmt,
		currentAccount.TableOrder,
		currentAccount.TableOrder+direction,
		currentAccount.TableOrder+direction,
		currentAccount.ID,
	)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
