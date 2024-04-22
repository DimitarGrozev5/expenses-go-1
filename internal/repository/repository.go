package repository

import (
	"database/sql"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

type DatabaseRepo interface {
	Close() error

	// User methods
	GetUserByEmail(email string) (models.User, error)
	Authenticate(email, testPassword string) (int, string, int, error)

	// Tags methods
	GetTags() ([]models.Tag, error)
	UpdateTags(tags []models.Tag, etx *sql.Tx) ([]models.Tag, error)

	// Expense methods
	GetExpenses() ([]models.Expense, error)
	AddExpense(expense models.Expense) error
	EditExpense(expense models.Expense) error
	DeleteExpense(id int) error
}
