package dbrepo

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/ctrlrepo"
)

type sqliteDBRepo struct {
	App *config.DBControllerConfig
	DB  *sql.DB
}

func NewSqliteRepo(app *config.DBControllerConfig, conn *sql.DB) ctrlrepo.ControllerRepository {
	return &sqliteDBRepo{
		App: app,
		DB:  conn,
	}
}

// Get user db path
func GetDBPath(dbPath, dbName string, fileOnly bool) string {

	if fileOnly {
		return fmt.Sprintf("%s%s.db", dbPath, dbName)
	}

	return fmt.Sprintf("%s%s.db?_fk=%s&_txlock=%s", dbPath, dbName, url.QueryEscape("true"), url.QueryEscape("exclusive"))
}
