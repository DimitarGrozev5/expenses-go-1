package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	cp "github.com/otiai10/copy"
)

var db *sql.DB

// Setup functions
func beforeAll() {
	if db == nil {
		// Get DB name
		dbName := "test"

		// Migrate db
		err := Migrate(dbName)
		if err != nil {
			log.Fatal(err)
		}

		// Create db
		db, err = sql.Open("sqlite3", fmt.Sprintf("%s.db?_fk=%s", dbName, url.QueryEscape("true")))
		if err != nil {
			log.Fatal(err)
		}

		// Copy db
		cp.Copy("test.db", "test_copy.db")
	}
}

func beforeEach() {
	// Get a copy of the db
	cp.Copy("test_copy.db", "test.db")
}

func afterAll() {
	// Close db conn
	db.Close()

	// Delete old DB
	os.Remove("test.db")
	os.Remove("test_copy.db")
}

// Setup multiple tests
func setupMultiple(t *testing.T, tests ...func(t *testing.T)) {
	beforeAll()
	defer afterAll()

	for _, test := range tests {
		beforeEach()
		test(t)
	}
}

func TestMigration(t *testing.T) {
	setupMultiple(
		t,

		// When you insert an account it should have the correct table_order
		func(t *testing.T) {
			// Insert account
			stmt := `INSERT INTO accounts (name) VALUES ('test account')`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert an account")
				return
			}

			var table_order int

			// Test account table order
			query := `SELECT table_order FROM accounts WHERE id=1`

			// Get rows
			row := db.QueryRow(query)

			// Scan row into model
			err = row.Scan(&table_order)
			if err != nil {
				t.Error("error getting account")
				return
			}

			// Test table_order
			if table_order != 1 {
				t.Errorf("table_order is %d; it has to be 1", table_order)
			}

			// Add two more accounts
			// Insert account
			stmt = `INSERT INTO accounts (name) VALUES ('test account1');
					INSERT INTO accounts (name) VALUES ('test account2');`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert two more accounts")
				return
			}

			// Get rows
			query = `SELECT table_order FROM accounts WHERE id > 1`

			// Get rows
			rows, err := db.Query(query)
			if err != nil {
				t.Error("couldn't get new accounts")
				return
			}

			var table_orders []int

			for rows.Next() {
				err = rows.Scan(&table_order)
				if err != nil {
					t.Error("couldn't scan new account")
					return
				}
				table_orders = append(table_orders, table_order)
			}

			if table_orders[0] != 2 || table_orders[1] != 3 {
				t.Errorf("table oreder of second and third account is %d, %d; has to be 2, 3", table_orders[0], table_orders[1])
				return
			}
		},

		// When you delete an account, the table order should update
		func(t *testing.T) {
			// Insert account
			stmt := `INSERT INTO accounts (name) VALUES ('test account1');
					 INSERT INTO accounts (name) VALUES ('test account2');
					 INSERT INTO accounts (name) VALUES ('test account3');
					 INSERT INTO accounts (name) VALUES ('test account4');`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert accounts;", err)
				return
			}

			// Delete second account
			stmt = `DELETE FROM accounts WHERE id=2`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete the second account")
				return
			}

			var table_orders []int

			// Get rows
			query := `SELECT table_order FROM accounts ORDER BY table_order ASC`

			// Get rows
			rows, err := db.Query(query)
			if err != nil {
				t.Error("couldn't get new accounts")
				return
			}

			for rows.Next() {
				var table_order int
				err = rows.Scan(&table_order)
				if err != nil {
					t.Error("couldn't scan new account")
					return
				}
				table_orders = append(table_orders, table_order)
			}

			// Check if order is ok
			if len(table_orders) > 3 {
				t.Errorf("too many items in table_orsers: %d, expected 3", len(table_orders))
				return
			}
			if table_orders[0] != 1 || table_orders[1] != 2 || table_orders[2] != 3 {
				t.Errorf("accounts order doesn't match; expected: 1, 2, 3; received: %d, %d, %d", table_orders[0], table_orders[1], table_orders[2])
				return
			}
		},

		// Can update account order
		func(t *testing.T) {
			// Insert account
			stmt := `INSERT INTO accounts (name) VALUES ('test account1');
					 INSERT INTO accounts (name) VALUES ('test account2');
					 INSERT INTO accounts (name) VALUES ('test account3');`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert accounts;", err)
				return
			}

			// Swap accounts order
			stmt = `UPDATE accounts SET table_order=1 WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't update accounts order;", err)
				return
			}

			var table_orders []int

			// Get rows
			query := `SELECT table_order FROM accounts ORDER BY id ASC`

			// Get rows
			rows, err := db.Query(query)
			if err != nil {
				t.Error("couldn't get accounts", err)
				return
			}

			for rows.Next() {
				var table_order int
				err = rows.Scan(&table_order)
				if err != nil {
					t.Error("couldn't scan new account")
					return
				}
				table_orders = append(table_orders, table_order)
			}

			// Check if order is ok
			if len(table_orders) > 3 {
				t.Errorf("too many items in table_orsers: %d, expected 3", len(table_orders))
				return
			}
			if table_orders[0] != 2 || table_orders[1] != 1 || table_orders[2] != 3 {
				t.Errorf("accounts order doesn't match; expected: 2, 1, 3; received: %d, %d, %d", table_orders[0], table_orders[1], table_orders[2])
				return
			}
		},

		// Don't allow account order to go bellow 1 and above the current max number
		func(t *testing.T) {
			// Insert account
			stmt := `INSERT INTO accounts (name) VALUES ('test account1');
					 INSERT INTO accounts (name) VALUES ('test account2');
					 INSERT INTO accounts (name) VALUES ('test account3');`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert accounts;", err)
				return
			}

			// Swap accounts order
			stmt = `UPDATE accounts SET table_order=0 WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("expected error when updating account table_order bellow 1")
				return
			}

			// Swap accounts order
			stmt = `UPDATE accounts SET table_order=4 WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("expected error when updating account table_order above current maximum")
				return
			}
		},
	)
}
