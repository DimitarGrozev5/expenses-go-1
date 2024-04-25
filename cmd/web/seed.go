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

		from_account	INTEGER		NOT NULL	REFERENCES accounts (id)
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

	// TODO: make default values of inital_amount 0; It's 100 for dev purposes
	// Define Category table
	stmt = `CREATE TABLE accounts (
		id				INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,
		
		name			TEXT		NOT NULL,
		initial_amount	NUMERIC		NOT NULL	DEFAULT 100	CHECK (initial_amount >= 0),
		current_amount	NUMERIC		NOT NULL	DEFAULT 100	CHECK (current_amount >= 0),

		usage_count		INTEGER		NOT NULL	DEFAULT 0,
		table_order		INTEGER		NOT NULL	DEFAULT -1,

		created_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}

	// Update accounts data when expenses change
	stmt = `	CREATE TRIGGER account_current_amount_add
					BEFORE INSERT
					ON expenses
				BEGIN
					UPDATE accounts SET
						current_amount = current_amount - new.amount,
						usage_count = usage_count + 1,
						updated_at = datetime('now')
					WHERE accounts.id = new.from_account;
				END;
				
				CREATE TRIGGER account_current_amount_remove
					BEFORE DELETE
					ON expenses
				BEGIN
					UPDATE accounts SET
						current_amount = current_amount + old.amount,
						usage_count = usage_count - 1,
						updated_at = datetime('now')
					WHERE accounts.id = old.from_account;
				END;
				
				CREATE TRIGGER account_current_amount_update_amount
					BEFORE UPDATE
					ON expenses
					WHEN
						old.amount <> new.amount AND
						old.from_account = new.from_account
				BEGIN
					UPDATE accounts SET
						current_amount = current_amount + old.amount - new.amount,
						updated_at = datetime('now')
					WHERE accounts.id = new.from_account;
				END;
				
				CREATE TRIGGER account_current_amount_update_account
					BEFORE UPDATE
					ON expenses
					WHEN old.from_account <> new.from_account
				BEGIN
					UPDATE accounts SET
						current_amount = current_amount + old.amount,
						usage_count = usage_count - 1
					WHERE accounts.id = old.from_account;
					
					UPDATE accounts SET
						current_amount = current_amount - new.amount,
						usage_count = usage_count + 1,
						updated_at = datetime('now')
					WHERE accounts.id = new.from_account;
				END;

				CREATE TRIGGER accounts_set_order_for_new_accounts
					AFTER INSERT
					ON accounts
				BEGIN
					UPDATE accounts SET
						table_order = (SELECT COUNT(*) FROM accounts),
						updated_at = datetime('now')
					WHERE accounts.table_order = -1;
				END;

				CREATE TRIGGER accounts_update_account_order
					BEFORE UPDATE
					ON accounts
				BEGIN
					SELECT
						CASE
							WHEN new.table_order < 1 THEN
								RAISE (ABORT, 'Cant move the first account up')
							WHEN new.table_order > (SELECT COUNT(*) from accounts) THEN
								RAISE (ABORT, 'Cant move the last account down')
						END;
				END;

				CREATE TRIGGER accounts_block_delete
					BEFORE DELETE
					ON accounts
				BEGIN
					SELECT
						CASE
							WHEN old.usage_count > 0 THEN
								RAISE (ABORT, 'Cant delete account that is being used')
						END;
				END;

				CREATE TRIGGER accounts_update_order_after_delete
					AFTER DELETE
					ON accounts
				BEGIN
					UPDATE accounts SET
						table_order = table_order - 1
					WHERE accounts.table_order > old.table_order;
				END;
	`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return
	}
}
