package repository

import (
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

type DatabaseRepo interface {
	Close() error

	// User methods
	GetUserByEmail(email string) (models.User, error)
	Authenticate(email, testPassword string) (int, string, error)

	// Expense methods
	GetExpenses() ([]models.Expense, error)
	AddExpense(expense models.Expense) error
}
