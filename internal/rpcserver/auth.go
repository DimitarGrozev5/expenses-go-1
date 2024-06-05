package rpcserver

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository/dbrepo"
	"github.com/golang-jwt/jwt/v5"
)

func (m *DatabaseServer) Authenticate(ctx context.Context, lc *models.LoginCredentials) (*models.LoginToken, error) {

	var loginResponse models.LoginToken

	// validate fields
	if len(lc.Email) == 0 || len(lc.Password) == 0 {
		return &loginResponse, fmt.Errorf("email and Password are required")
	}

	// Check if user DB exists
	_, err := os.Stat(dbrepo.GetUserDBPath(m.App.DBPath, lc.Email, true))
	if errors.Is(err, os.ErrNotExist) {

		// Write to error log
		m.App.ErrorLog.Println(err)

		// Return error
		return &loginResponse, fmt.Errorf("invalid login credentials")
	}

	// Create user connection
	dbconn, err := driver.ConnectSQL(dbrepo.GetUserDBPath(m.App.DBPath, lc.Email, false))
	if err != nil {

		// Write to error log
		m.App.ErrorLog.Println(err)

		// Return error
		return &loginResponse, fmt.Errorf("server error")
	}

	// Get db repo
	repo := dbrepo.NewSqliteRepo(m.App, lc.Email, dbconn.SQL)

	// Authenticate user
	_, _, dbVersion, err := repo.Authenticate(lc.Password)
	if err != nil {

		// Write to error log
		m.App.ErrorLog.Println(err)

		return &loginResponse, fmt.Errorf("invalid login credentials")
	}

	// Get user key
	key := dbrepo.GetUserKey(lc.Email)

	// Add connection to repo
	m.App.DBConnections[key] = repo

	// Crate JWT to authenticate user
	t := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"userKey":   key,
			"dbVersion": dbVersion,
		})
	jwt, err := t.SignedString(m.App.JWTKey)
	if err != nil {
		return &loginResponse, err
	}

	// Add token to response
	loginResponse.Token = jwt

	return &loginResponse, nil
}
