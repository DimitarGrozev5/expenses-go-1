package models

import (
	"database/sql"
	"time"
)

// NewExpense model
type NewExpense struct {
	Amount float64
	Tag    string
	Date   string
}

// Users
type User struct {
	ID        int
	Email     string
	Password  string
	DBVersion int
	FreeFunds float64
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// Expenses
type Expense struct {
	ID             int
	Amount         float64
	Date           time.Time
	Tags           []Tag
	FromAccountId  int
	FromAccount    Account
	FromCategoryId int
	FromCategory   Category
	CreatedAt      time.Time
	UpdatedAt      sql.NullTime
}

// Tags
type Tag struct {
	ID         int
	Name       string
	UsageCount int
	CreatedAt  time.Time
	UpdatedAt  sql.NullTime
}

// Expense-Tag relation
type ExpenseToTagRealtion struct {
	ID        int
	ExpenseId int
	TagId     int
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// Accounts
type Account struct {
	ID int

	Name          string
	CurrentAmount float64

	UsageCount int
	TableOrder int

	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// Categories
type Category struct {
	ID int

	Name string

	BudgetInput   float64
	LastInputDate time.Time
	InputInterval time.Duration

	SpendingLimit float64
	SpendingLeft  float64

	InitialAmount float64
	CurrentAmount float64

	TableOrder int
	CreatedAt  time.Time
	UpdatedAt  sql.NullTime
}

type CategoryOverview struct {
	ID   int
	Name string

	BudgetInput   float64
	InputInterval int
	InputPriodId  int

	SpendingLimit float64
	SpendingLeft  float64

	PeriodStart time.Time
	PeriodEnd   time.Time

	InitialAmount float64
	CurrentAmount float64

	CanBeDeleted bool
	TableOrder   int
}

// Time periods
type TimePeriod struct {
	ID        int
	Period    string
	Caption   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
