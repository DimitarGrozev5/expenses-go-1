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
	query := `SELECT id, name, current_amount FROM accounts`

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

		err = rows.Scan(&account.ID, &account.Name, &account.CurrentAmount)
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

func (m *sqliteDBRepo) AddAccount(account models.Account) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query to insert account
	stmt := `INSERT INTO accounts(name, initial_amount, current_amount) VALUES($1, $2, $3)`

	// Execute query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		account.Name,
		account.InitialAmount,
		account.InitialAmount,
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
	return nil
}

func (m *sqliteDBRepo) TransferFunds(fromAccount, toAccount models.Account, amount float64) error {
	return nil
}
