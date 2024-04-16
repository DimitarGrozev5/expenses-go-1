package models

import "time"

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
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Expenses
type Expense struct {
	ID        int
	Amount    float64
	Tag       string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	// CategoryID int
	// Category   Category
	// AccountID  int
	// Account    Amount
}
