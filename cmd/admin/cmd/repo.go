package cmd

import (
	"github.com/dimitargrozev5/expenses-go-1/internal/ctrlrepo"
	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
)

type CommandRepository struct {
	CtrlDBConn *driver.DB
	CtrlDB     ctrlrepo.ControllerRepository
}

var Repo CommandRepository

func InitCmdRepo(ctrlDBConn *driver.DB, ctrldb ctrlrepo.ControllerRepository) {
	Repo.CtrlDBConn = ctrlDBConn
	Repo.CtrlDB = ctrldb
}
