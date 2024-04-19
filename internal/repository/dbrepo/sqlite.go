package dbrepo

import (
	"context"
	"errors"
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

// Get expenses ordered by date
func (m *sqliteDBRepo) GetExpenses() ([]models.Expense, error) {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	query := `SELECT id, amount, label, date FROM expenses ORDER BY date DESC`

	// Get rows
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define exenses slice
	expenses := make([]models.Expense, 0)

	// Scan rows
	for rows.Next() {
		var expense models.Expense

		err = rows.Scan(&expense.ID, &expense.Amount, &expense.Label, &expense.Date)
		if err != nil {
			log.Fatal(err)
		}

		// Add to slice
		expenses = append(expenses, expense)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

// Add expense
func (m *sqliteDBRepo) AddExpense(expense models.Expense) error {
	// Define context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Define query
	stmt := `INSERT INTO expenses(amount, label, date) VALUES($1, $2, $3)`

	// Execute query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		expense.Amount,
		expense.Label,
		expense.Date,
	)

	if err != nil {
		return err
	}

	return nil
}
