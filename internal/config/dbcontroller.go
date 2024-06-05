package config

import (
	"crypto/ecdsa"
	"log"

	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
)

// AppConfig holds the application config
type DBControllerConfig struct {
	InProduction  bool
	DBPath        string
	JWTKey        *ecdsa.PrivateKey
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	DBConnections map[string]repository.DatabaseRepo
}
