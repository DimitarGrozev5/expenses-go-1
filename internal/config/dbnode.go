package config

import (
	"log"

	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
	"github.com/dimitargrozev5/expenses-go-1/internal/repository"
)

// AppConfig holds the application config
type DBNodeConfig struct {
	NodeID            int64
	InProduction      bool
	ControllerAddress string
	DBPath            string
	JWTSecretKey      []byte //*ecdsa.PrivateKey
	InfoLog           *log.Logger
	ErrorLog          *log.Logger
	DBConnections     map[string]*driver.DB
	DBRepos           map[string]repository.DatabaseRepo
}

func (c DBNodeConfig) GetJWTSecretKey() []byte {
	return c.JWTSecretKey
}

func (c DBNodeConfig) GetInProduction() bool {
	return c.InProduction
}
