package config

import (
	"log"

	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
)

// AppConfig holds the application config
type DBControllerConfig struct {
	InProduction  bool
	DBPath        string
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	DBConnections map[string]repository.DatabaseRepo
}
