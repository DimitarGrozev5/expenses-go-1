package config

import (
	"database/sql"
	"log"

	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
)

// AppConfig holds the application config
type DBControllerConfig struct {
	InProduction   bool
	DBPath         string
	DBName         string
	CtrlDB         *sql.DB
	MigrationsPath string
	JWTSecretKey   []byte //*ecdsa.PrivateKey
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	DBConnections  map[string]*driver.DB
	DBRepos        map[string]repository.DatabaseRepo
}

func (c DBControllerConfig) GetJWTSecretKey() []byte {
	return c.JWTSecretKey
}
