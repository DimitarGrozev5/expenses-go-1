package grpcrepository

import (
	"context"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

type DatabaseRepo interface {
	// User methods
	GetUser(ctx context.Context, empty models.GrpcEmpty) (models.GrpcUser, error)
	Authenticate(ctx context.Context, loginCredentials models.LoginCredentials) (models.LoginToken, error)
	ModifyFreeFunds(ctx context.Context, params models.ModifyFreeFundsParams) (models.GrpcEmpty, error)

	// Tags methods
	GetTags(ctx context.Context, empty models.GrpcEmpty) (models.GetTagsReturns, error)

	// Expense methods
	GetExpenses(ctx context.Context, param models.GrpcEmpty) models.GetExpensesReturns
	AddExpense(ctx context.Context, param models.ExpensesParams) models.GrpcEmpty
	EditExpense(ctx context.Context, param models.ExpensesParams) models.GrpcEmpty
	DeleteExpense(ctx context.Context, param models.DeleteExpenseParams) models.GrpcEmpty

	// Account methods
	GetAccounts(ctx context.Context, params models.GetAccountsParams) models.GetAccountsReturns
	AddAccount(ctx context.Context, params models.AddAccountParams) models.GrpcEmpty
	EditAccountName(ctx context.Context, params models.EditAccountNameParams) models.GrpcEmpty
	DeleteAccount(ctx context.Context, params models.DeleteAccountParams) models.GrpcEmpty
	TransferFunds(ctx context.Context, params models.TransferFundsParams) models.GrpcEmpty
	ReorderAccount(ctx context.Context, params models.ReorderAccountParams) models.GrpcEmpty

	// Categories methods
	GetCategoriesCount(ctx context.Context, params models.GrpcEmpty) models.GetCategoriesCountReturns
	GetCategories(ctx context.Context, params models.GrpcEmpty) models.GetCategoriesReturns
	GetCategoriesOverview(ctx context.Context, params models.GrpcEmpty) models.GetCategoriesOverviewReturns
	AddCategory(ctx context.Context, params models.AddCategoryParams) models.GrpcEmpty
	ReorderCategory(ctx context.Context, params models.ReorderCategoryParams) models.GrpcEmpty
	DeleteCategory(ctx context.Context, params models.DeleteCategoryParams) models.GrpcEmpty
	ResetCategories(ctx context.Context, params models.ResetCategoriesParams) models.GrpcEmpty

	// Time periods
	GetTimePeriods(ctx context.Context, empty models.GrpcEmpty) models.GetTimePeriodsReturns
}
