package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

func Seed(DBPath string) {

	// Get DB name
	dbName := DBPath + "asd@asd.asd"

	// Migrate db
	err := Migrate(dbName)

	// Create db
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s.db?_fk=%s", dbName, url.QueryEscape("true")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Seed user
	stmt := `INSERT INTO user (email, password, db_version) VALUES ($1, $2, $3)`

	// Get password
	password, _ := bcrypt.GenerateFromPassword([]byte("asd"), bcrypt.DefaultCost)

	// Execute query
	_, err = db.Exec(stmt, "asd@asd.asd", password, 0)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}

}
