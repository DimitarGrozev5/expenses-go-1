package repository

import (
	"database/sql"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

type DatabaseRepo interface {
	Close() error

	// User methods
	GetUser() (models.User, error)
	Authenticate(testPassword string) (int, string, int, error)
	ModifyFreeFunds(amount float64, toAccountId int, tagName string) error

	// Tags methods
	GetTags() ([]models.Tag, error)
	UpdateTags(tags []string, etx *sql.Tx) ([]models.Tag, error)

	// Expense methods
	GetExpenses() ([]models.Expense, error)
	AddExpense(expense models.Expense, tags []string) error
	EditExpense(expense models.Expense, tags []string) error
	DeleteExpense(id int) error

	// Account methods
	GetAccounts(orderByPopularity bool) ([]models.Account, error)
	AddAccount(name string) error
	EditAccountName(id int, name string) error
	DeleteAccount(id int) error
	TransferFunds(fromAccount, toAccount models.Account, amount float64) error
	ReorderAccount(currentAccount models.Account, direction int) error

	// Categories methods
	GetCategoriesCount() (int, error)
	GetCategories() ([]models.Category, error)
	GetCategoriesOverview() ([]models.CategoryOverview, error)
	AddCategory(name string, budgetInput float64, spendingLimit float64, inputInterval int, inputPeriod int) error
	ReorderCategory(categoryid int, new_order int) error
	DeleteCategory(id int) error
	ResetCategory(amount float64, categoryId int, budgetInput float64, inputInterval int, inputPeriod int, spendingLimit float64, etx *sql.Tx) error
	ResetCategories(cateogries []models.ResetCategoryData) error

	// Time periods
	GetTimePeriods() ([]models.TimePeriod, error)
}
