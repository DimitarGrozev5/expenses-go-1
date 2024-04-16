package dbrepo

import (
	"database/sql"

	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
)

type sqliteDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewSqliteRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &sqliteDBRepo{
		App: a,
		DB:  conn,
	}
}
