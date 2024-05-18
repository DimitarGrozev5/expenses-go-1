package dbrepo

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
)

type sqliteDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewSqliteRepo(app *config.AppConfig, user string, conn *sql.DB) repository.DatabaseRepo {
	return &sqliteDBRepo{
		App: app,
		DB:  conn,
	}
}

// Get user key
func GetUserKey(user string) string {
	return user
}

// Get user db path
func GetUserDBPath(dbPath, user string) string {
	// Get user key
	userKey := GetUserKey(user)

	return fmt.Sprintf("%s%s.db?_fk=%s&_txlock=%s", dbPath, userKey, url.QueryEscape("true"), url.QueryEscape("exclusive"))
}
