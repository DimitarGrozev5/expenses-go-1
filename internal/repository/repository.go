package repository

import (
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

type DatabaseRepo interface {
	// User methods
	GetUser(empty *models.GrpcEmpty) (*models.GrpcUser, error)
	Authenticate(testPassword string) (int64, string, int64, error)
	ModifyFreeFunds(params *models.ModifyFreeFundsParams) (*models.GrpcEmpty, error)

	// Tags methods
	GetTags(empty *models.GrpcEmpty) (*models.GetTagsReturns, error)

	// Expense methods
	GetExpenses(param *models.GrpcEmpty) (*models.GetExpensesReturns, error)
	AddExpense(param *models.ExpensesParams) (*models.GrpcEmpty, error)
	EditExpense(param *models.ExpensesParams) (*models.GrpcEmpty, error)
	DeleteExpense(param *models.DeleteExpenseParams) (*models.GrpcEmpty, error)

	// Account methods
	GetAccounts(params *models.GetAccountsParams) (*models.GetAccountsReturns, error)
	AddAccount(params *models.AddAccountParams) (*models.GrpcEmpty, error)
	EditAccountName(params *models.EditAccountNameParams) (*models.GrpcEmpty, error)
	DeleteAccount(params *models.DeleteAccountParams) (*models.GrpcEmpty, error)
	TransferFunds(params *models.TransferFundsParams) (*models.GrpcEmpty, error)
	ReorderAccount(params *models.ReorderAccountParams) (*models.GrpcEmpty, error)

	// Categories methods
	// GetCategoriesCount(params *models.GrpcEmpty) (*models.GetCategoriesCountReturns, error)
	GetCategories(params *models.GrpcEmpty) (*models.GetCategoriesReturns, error)
	GetCategoriesOverview(params *models.GrpcEmpty) (*models.GetCategoriesOverviewReturns, error)
	AddCategory(params *models.AddCategoryParams) (*models.GrpcEmpty, error)
	ReorderCategory(params *models.ReorderCategoryParams) (*models.GrpcEmpty, error)
	DeleteCategory(params *models.DeleteCategoryParams) (*models.GrpcEmpty, error)
	ResetCategories(params *models.ResetCategoriesParams) (*models.GrpcEmpty, error)

	// Time periods
	GetTimePeriods(empty *models.GrpcEmpty) (*models.GetTimePeriodsReturns, error)
}
