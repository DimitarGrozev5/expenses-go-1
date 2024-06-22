package cmd

import "github.com/dimitargrozev5/expenses-go-1/internal/ctrlrepo"

type CommandRepository struct {
	CtrlDB ctrlrepo.ControllerRepository
}

var Repo CommandRepository

func InitCmdRepo(ctrldb ctrlrepo.ControllerRepository) {
	Repo.CtrlDB = ctrldb
}
