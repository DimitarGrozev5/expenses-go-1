package main

import (
	"flag"
	"log"

	"github.com/dimitargrozev5/expenses-go-1/cmd/admin/cmd"
	"github.com/dimitargrozev5/expenses-go-1/internal/config"
	"github.com/dimitargrozev5/expenses-go-1/internal/ctrlrepo/dbrepo"
	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
)

var dbPath = flag.String("db-path", "./db/", "Path to folder containing sqlite databases")
var dbCtrlName = flag.String("db-name", "ctrl.db", "Controller DB name")

func main() {
	// Connect to controller db
	db, err := driver.ConnectSQL(dbrepo.GetDBPath(*dbPath, *dbCtrlName, false))
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	// Start Controller repo
	repo := dbrepo.NewSqliteRepo(&config.DBControllerConfig{}, db.SQL)

	// Init commands
	cmd.InitCmdRepo(repo)

	cmd.Execute()
}
