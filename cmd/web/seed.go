package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func Seed(DBPath string) {

	// Delete old DB
	os.Remove(DBPath + "asd@asd.asd.db")

	_, err := os.OpenFile(DBPath+"asd@asd.asd.db", os.O_RDONLY, os.ModeType)
	if err == nil {
		log.Printf("DB not deleted!!!!!!!!!!!")
		return
	}

	// Create db
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s%s.db?_fk=%s", DBPath, "asd@asd.asd", url.QueryEscape("true")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create user table
	stmt := `CREATE TABLE user (
					id			INTEGER					NOT NULL	PRIMARY KEY		AUTOINCREMENT,

					email		TEXT		UNIQUE		NOT NULL,
					password	TEXT					NOT NULL,
					db_version	INTEGER

					created_at	DATETIME				NOT NULL	DEFAULT CURRENT_TIMESTAMP,
					updated_at	DATETIME				NOT NULL	DEFAULT CURRENT_TIMESTAMP
				)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}

	// Insert in user table
	stmt = `INSERT INTO user (email, password, db_version) VALUES ($1, $2, $3)`

	// Get password
	password, _ := bcrypt.GenerateFromPassword([]byte("asd"), bcrypt.DefaultCost)

	// Execute query
	_, err = db.Exec(stmt, "asd@asd.asd", password, 0)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}

	// Create table expenses
	stmt = `CREATE TABLE expenses (
		id				INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

		amount			NUMERIC		NOT NULL,
		date			DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,

		from_account	INTEGER		NOT NULL	REFERENCES expenses (id)
													ON DELETE RESTRICT,
		created_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}

	// Create table tags
	stmt = `CREATE TABLE tags (
		id			INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

		name		TEXT		NOT NULL	UNIQUE,
		usage_count	INTEGER		NOT NULL					DEFAULT 0,
		last_used	DATETIME	NOT NULL					DEFAULT CURRENT_TIMESTAMP,

		created_at	DATETIME	NOT NULL					DEFAULT CURRENT_TIMESTAMP,
		updated_at	DATETIME	NOT NULL					DEFAULT CURRENT_TIMESTAMP
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}

	// Create table expense tags
	stmt = `CREATE TABLE expense_tags (
		id			INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

		expense_id	INTEGER		NOT NULL	REFERENCES expenses (id)
												ON DELETE CASCADE
												ON UPDATE CASCADE,
		tag_id		INTEGER		NOT NULL	REFERENCES tags (id)
												ON DELETE CASCADE
												ON UPDATE CASCADE,

		created_at	DATETIME	NOT NULL					DEFAULT CURRENT_TIMESTAMP,
		updated_at	DATETIME	NOT NULL					DEFAULT CURRENT_TIMESTAMP
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}

	// Update tag usage_count with trigger
	stmt = `	CREATE TRIGGER tag_usage_count_insert
					AFTER INSERT
					ON expense_tags
				BEGIN
					UPDATE tags SET
						usage_count = usage_count + 1,
						last_used = datetime('now'),
						updated_at = datetime('now')
					WHERE tags.id = new.tag_id;
				END;
				
				CREATE TRIGGER tag_usage_count_update
					AFTER UPDATE
					ON expense_tags
					WHEN old.tag_id <> new.tag_id
				BEGIN
					UPDATE tags SET
						usage_count = usage_count + 1,
						last_used = datetime('now'),
						updated_at = datetime('now')
					WHERE tags.id = new.tag_id;
					
					UPDATE tags SET
						usage_count = usage_count - 1,
						updated_at = datetime('now')
					WHERE tags.id = old.tag_id;
				END;
				
				CREATE TRIGGER tag_usage_count_delete
					AFTER DELETE
					ON expense_tags
				BEGIN
					UPDATE tags SET
						usage_count = usage_count - 1,
						updated_at = datetime('now')
					WHERE tags.id = old.tag_id;
				END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}

	// Define Category table
	stmt = `CREATE TABLE accounts (
		id				INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,
		
		name			TEXT		NOT NULL,
		initial_amount	NUMERIC		NOT NULL,

		created_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}
}
