package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
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
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s.db?_fk=%s", app.DBPath+app.DBName, url.QueryEscape("true")))
	if err != nil {
		log.Fatal(err)
	}

	// Add connection to state
	app.CtrlDB = db

	// Setup new context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get DB Version
	row := db.QueryRowContext(ctx, `PRAGMA user_version`)

	// Pull data from row
	var userVersion int
	err = row.Scan(&userVersion)
	if err != nil {
		log.Fatal(err)
	}

	// Store migrations
	migrations := make([]string, 0)

	// Get all migrations from current db version up
	i := userVersion
	for {

		// Update count
		i++

		// Try to get file
		file, err := os.ReadFile(app.MigrationsPath + fmt.Sprintf("ctrl-up-%d.sql", i))

		// Exit if migration not found
		if os.IsNotExist(err) {
			break
		}

		// Panic if other error
		if err != nil {
			log.Fatal(err)
		}

		// Add file to migration
		migrations = append(migrations, string(file))
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// Run migrations from current version up
	for _, migration := range migrations {

		// Run query
		_, err = tx.ExecContext(ctx, migration)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Print migration message
	app.InfoLog.Printf("Migrations performed: %d; Controller DB Version: %d", len(migrations), i-1)

	// Commit migrations
	tx.Commit()
}
