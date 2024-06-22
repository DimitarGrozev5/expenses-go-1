package main

import (
	"log"
	"os"

	"github.com/dimitargrozev5/expenses-go-1/internal/ctrlrepo/dbrepo"
	"github.com/dimitargrozev5/expenses-go-1/internal/driver"
)

func setupDb() {
	// Check for db folder
	_, err := os.Stat(app.DBPath)
	if os.IsNotExist(err) {
		// Create db folder
		if err := os.Mkdir(app.DBPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// Check if db file exists
	// _, err = os.Stat(app.DBPath + app.DBName)
	// // if os.IsNotExist(err) {
	// // 	log.Fatal("File does not exist.")
	// // }

	// Open db connection
	db, err := driver.ConnectSQL(dbrepo.GetDBPath(app.DBPath, app.DBName, false))
	if err != nil {
		log.Fatal(err)
	}

	// Add connection to state
	app.CtrlDB = db.SQL

	// Start Controller repo
	repo := dbrepo.NewSqliteRepo(&app, db.SQL)

	// Add repo to state
	app.CtrlDBRepo = repo
}
