package dbrepo

import (
	"context"
	"errors"
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
	query := `SELECT id, email, password, db_version
				FROM user WHERE email = $1`

	// Get row
	row := m.DB.QueryRowContext(ctx, query, email)

	// Scan row into model
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.DBVersion,
	)

	// Check for error
	if err != nil {
		return u, err
	}

	// Return user
	return u, nil
}

// Authenticate user
func (m *sqliteDBRepo) Authenticate(email, testPassword string) (int, string, int, error) {
	// Get user
	u, err := m.GetUserByEmail(email)
	if err != nil {
		return 0, "", 0, err
	}

	// Set variable
	id := u.ID
	hashedPassword := u.Password
	dbVersion := u.DBVersion

	// Check if password matches
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, "", 0, err
	}

	return id, hashedPassword, dbVersion, nil
}
