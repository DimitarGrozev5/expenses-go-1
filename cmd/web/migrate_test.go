package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
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
	defer func() {
		if r := recover(); r != nil {
			afterAll()
		}
		afterAll()
	}()

	for _, test := range tests {
		beforeEach()
		test(t)
	}
}

func TestMigration(t *testing.T) {
	setupMultiple(
		t,

		/*
		 * Account tests
		 *
		 */
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

		/*
		 * Categories tests
		 *
		 */
		// When you insert a category it should have the correct table_order
		func(t *testing.T) {
			// Insert account
			stmt := `INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category', 100, 1, 2, 100, 1, 2)`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert a category", err)
				return
			}

			var table_order int

			// Test category table order
			query := `SELECT table_order FROM categories WHERE id=1`

			// Get rows
			row := db.QueryRow(query)

			// Scan row into model
			err = row.Scan(&table_order)
			if err != nil {
				t.Error("error getting category")
				return
			}

			// Test table_order
			if table_order != 1 {
				t.Errorf("table_order is %d; it has to be 1", table_order)
			}

			// Add two more categories
			stmt = `INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category1', 100, 1, 2, 100, 1, 2);
					INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category2', 100, 1, 2, 100, 1, 2);`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert two more categories", err)
				return
			}

			// Get rows
			query = `SELECT table_order FROM categories WHERE id > 1`

			// Get rows
			rows, err := db.Query(query)
			if err != nil {
				t.Error("couldn't get new category", err)
				return
			}

			var table_orders []int

			for rows.Next() {
				err = rows.Scan(&table_order)
				if err != nil {
					t.Error("couldn't scan new categories", err)
					return
				}
				table_orders = append(table_orders, table_order)
			}

			if table_orders[0] != 2 || table_orders[1] != 3 {
				t.Errorf("table oreder of second and third category is %d, %d; has to be 2, 3", table_orders[0], table_orders[1])
				return
			}
		},

		// When you delete a category, the table order should update
		func(t *testing.T) {
			// Insert account
			stmt := `INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category1', 100, 1, 2, 100, 1, 2);
					 INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category2', 100, 1, 2, 100, 1, 2);
					 INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category3', 100, 1, 2, 100, 1, 2);
					 INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category4', 100, 1, 2, 100, 1, 2);`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert categories;", err)
				return
			}

			// Delete second account
			stmt = `DELETE FROM categories WHERE id=2`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete the second category", err)
				return
			}

			var table_orders []int

			// Get rows
			query := `SELECT table_order FROM categories ORDER BY id ASC`

			// Get rows
			rows, err := db.Query(query)
			if err != nil {
				t.Error("couldn't get new categories")
				return
			}

			for rows.Next() {
				var table_order int
				err = rows.Scan(&table_order)
				if err != nil {
					t.Error("couldn't scan new categories")
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
				t.Errorf("categories order doesn't match; expected: 1, 2, 3; received: %d, %d, %d", table_orders[0], table_orders[1], table_orders[2])
				return
			}
		},

		// Can update category order
		func(t *testing.T) {
			// Insert account
			stmt := `INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category1', 100, 1, 2, 100, 1, 2);
					 INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category2', 100, 1, 2, 100, 1, 2);
					 INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category3', 100, 1, 2, 100, 1, 2);`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert categories;", err)
				return
			}

			// Swap categories order
			stmt = `UPDATE categories SET table_order=1 WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't update categories order;", err)
				return
			}

			var table_orders []int

			// Get rows
			query := `SELECT table_order FROM categories ORDER BY id ASC`

			// Get rows
			rows, err := db.Query(query)
			if err != nil {
				t.Error("couldn't get categories", err)
				return
			}

			for rows.Next() {
				var table_order int
				err = rows.Scan(&table_order)
				if err != nil {
					t.Error("couldn't scan new category")
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
				t.Errorf("caegories order doesn't match; expected: 2, 1, 3; received: %d, %d, %d", table_orders[0], table_orders[1], table_orders[2])
				return
			}
		},

		// Don't allow category order to go bellow 1 and above the current max number
		func(t *testing.T) {
			// Insert category
			stmt := `INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category1', 100, 1, 2, 100, 1, 2);
					 INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category2', 100, 1, 2, 100, 1, 2);
					 INSERT INTO categories ( name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category3', 100, 1, 2, 100, 1, 2);`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert categories;", err)
				return
			}

			// Swap categories order
			stmt = `UPDATE categories SET table_order=0 WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("expected error when updating category table_order bellow 1")
				return
			}

			// Swap category order
			stmt = `UPDATE category SET table_order=4 WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("expected error when updating category table_order above current maximum")
				return
			}
		},

		/*
		 * Test Expense interactions
		 *
		 */
		// When you add an Expense with tag relations it changes the related tags usage_count value
		// You can't delete tags, that are used
		func(t *testing.T) {
			// Run init
			err := beforeExpenseTest(t)
			if err != nil {
				t.Error(err)
				return
			}

			// Get current tag values
			query := `SELECT name, usage_count FROM tags`

			// Get tags
			var initialTags []models.Tag

			// Get rows
			rows, err := db.Query(query)
			if err != nil {
				t.Error("couldn't get tags", err)
				return
			}

			for rows.Next() {
				var tag models.Tag
				err = rows.Scan(&tag.Name, &tag.UsageCount)
				if err != nil {
					t.Error("couldn't scan tag", err)
					return
				}
				initialTags = append(initialTags, tag)
			}

			// Check tags initial usage_count
			if len(initialTags) > 2 {
				t.Errorf("too many tags: %d; expected 2", len(initialTags))
				return
			}
			if initialTags[0].UsageCount != 0 || initialTags[1].UsageCount != 0 {
				t.Errorf("unexpected tags initial usage_count; recevied: %d, %d; expected 0, 0", initialTags[0].UsageCount, initialTags[1].UsageCount)
			}

			// Add expense with relation to tag 1
			stmt := `INSERT INTO expenses (amount, from_account, from_category) VALUES (10, 1, 1);
					 INSERT INTO expense_tags (expense_id, tag_id) VALUES (1, 1);`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert expense and expense_tags;", err)
				return
			}

			// Get tags
			var tags1 []models.Tag

			// Get rows
			rows, err = db.Query(query)
			if err != nil {
				t.Error("couldn't get tags", err)
				return
			}

			for rows.Next() {
				var tag models.Tag
				err = rows.Scan(&tag.Name, &tag.UsageCount)
				if err != nil {
					t.Error("couldn't scan tag", err)
					return
				}
				tags1 = append(tags1, tag)
			}

			// Check tags initial usage_count
			if len(tags1) > 2 {
				t.Errorf("too many tags: %d; expected 2", len(tags1))
				return
			}
			if tags1[0].UsageCount != 1 || tags1[1].UsageCount != 0 {
				t.Errorf("unexpected tags initial usage_count; recevied: %d, %d; expected 1, 0", tags1[0].UsageCount, tags1[1].UsageCount)
			}

			// Add expense with relation to tag 1 and 2
			stmt = `INSERT INTO expenses (amount, from_account, from_category) VALUES (10, 1, 1);
					 INSERT INTO expense_tags (expense_id, tag_id) VALUES (2, 1);
					 INSERT INTO expense_tags (expense_id, tag_id) VALUES (2, 2);`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert expense and expense_tags;", err)
				return
			}

			// Get tags
			var tags2 []models.Tag

			// Get rows
			rows, err = db.Query(query)
			if err != nil {
				t.Error("couldn't get tags", err)
				return
			}

			for rows.Next() {
				var tag models.Tag
				err = rows.Scan(&tag.Name, &tag.UsageCount)
				if err != nil {
					t.Error("couldn't scan tag", err)
					return
				}
				tags2 = append(tags2, tag)
			}

			// Check tags initial usage_count
			if len(tags2) > 2 {
				t.Errorf("too many tags: %d; expected 2", len(tags2))
				return
			}
			if tags2[0].UsageCount != 2 || tags2[1].UsageCount != 1 {
				t.Errorf("unexpected tags initial usage_count; recevied: %d, %d; expected 2, 1", tags2[0].UsageCount, tags2[1].UsageCount)
			}

			// Delete from expense tags
			stmt = `DELETE FROM expense_tags WHERE id=3`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete expense tags;", err)
				return
			}

			// Get tags
			var tags3 []models.Tag

			// Get rows
			rows, err = db.Query(query)
			if err != nil {
				t.Error("couldn't get tags", err)
				return
			}

			for rows.Next() {
				var tag models.Tag
				err = rows.Scan(&tag.Name, &tag.UsageCount)
				if err != nil {
					t.Error("couldn't scan tag", err)
					return
				}
				tags3 = append(tags3, tag)
			}

			// Check tags initial usage_count
			if len(tags3) > 2 {
				t.Errorf("too many tags: %d; expected 2", len(tags3))
				return
			}
			if tags3[0].UsageCount != 2 || tags3[1].UsageCount != 0 {
				t.Errorf("unexpected tags initial usage_count; recevied: %d, %d; expected 2, 0", tags3[0].UsageCount, tags3[1].UsageCount)
				return
			}

			// Remove expense
			stmt = `DELETE FROM expenses WHERE id=1`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete expense;", err)
				return
			}

			// Get tags
			var tags4 []models.Tag

			// Get rows
			rows, err = db.Query(query)
			if err != nil {
				t.Error("couldn't get tags", err)
				return
			}

			for rows.Next() {
				var tag models.Tag
				err = rows.Scan(&tag.Name, &tag.UsageCount)
				if err != nil {
					t.Error("couldn't scan tag", err)
					return
				}
				tags4 = append(tags4, tag)
			}

			// Check tags initial usage_count
			if len(tags4) > 2 {
				t.Errorf("too many tags: %d; expected 2", len(tags4))
				return
			}
			if tags4[0].UsageCount != 1 || tags4[1].UsageCount != 0 {
				t.Errorf("unexpected tags initial usage_count; recevied: %d, %d; expected 1, 0", tags4[0].UsageCount, tags4[1].UsageCount)
				return
			}

			// Delete second tag
			stmt = `DELETE FROM tags WHERE id=2`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete tag;", err)
				return
			}

			// Get tags
			var tags5 []models.Tag

			// Get rows
			rows, err = db.Query(query)
			if err != nil {
				t.Error("couldn't get tags", err)
				return
			}

			for rows.Next() {
				var tag models.Tag
				err = rows.Scan(&tag.Name, &tag.UsageCount)
				if err != nil {
					t.Error("couldn't scan tag", err)
					return
				}
				tags5 = append(tags5, tag)
			}

			// Check tags initial usage_count
			if len(tags5) > 1 {
				t.Errorf("too many tags: %d; expected 1", len(tags5))
				return
			}

			// Delete first tag
			stmt = `DELETE FROM tags WHERE id=1`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("shouldn't be able to delete tag that has usage_count above zero", err)
				return
			}
		},

		// When you add an Expense, the amount gets reflected in the related account current_amount and usage_count and in the category current_amount
		func(t *testing.T) {
			// Run init
			err := beforeExpenseTest(t)
			if err != nil {
				t.Error(err)
				return
			}

			// Add expense
			stmt := `INSERT INTO expenses (amount, from_account, from_category) VALUES (10, 1, 1);`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert expense;", err)
				return
			}

			getAccounts := func() ([]models.Account, error) {

				// Get accounts
				query := `SELECT current_amount FROM accounts WHERE id=1;`

				// Set variable
				var accounts []models.Account

				// Get rows
				rows, err := db.Query(query)
				if err != nil {
					t.Error("couldn't get account", err)
					return nil, err
				}

				for rows.Next() {
					var account models.Account
					err = rows.Scan(&account.CurrentAmount)
					if err != nil {
						t.Error("couldn't scan account", err)
						return nil, err
					}
					accounts = append(accounts, account)
				}

				return accounts, nil
			}

			accounts, err := getAccounts()
			if err != nil {
				return
			}

			// Check if account current amount is changed
			if accounts[0].CurrentAmount != 90 {
				t.Errorf("account current amount is wrong; expected 90; received: %f", accounts[0].CurrentAmount)
				return
			}

			getCats := func() ([]models.Category, error) {
				// Get categories
				query := `SELECT current_amount FROM categories WHERE id=1;`

				// Set variable
				var categories []models.Category

				// Get rows
				rows, err := db.Query(query)
				if err != nil {
					t.Error("couldn't get categories", err)
					return nil, err
				}

				for rows.Next() {
					var category models.Category
					err = rows.Scan(&category.CurrentAmount)
					if err != nil {
						t.Error("couldn't scan category", err)
						return nil, err
					}
					categories = append(categories, category)
				}

				return categories, nil
			}

			categories, err := getCats()
			if err != nil {
				return
			}

			// Check if account current amount is changed
			if categories[0].CurrentAmount != 90 {
				t.Errorf("category current amount is wrong; expected 90; received: %f", categories[0].CurrentAmount)
				return
			}

			// Add expenses
			stmt = `INSERT INTO expenses (amount, from_account, from_category) VALUES (20, 1, 1);
					INSERT INTO expenses (amount, from_account, from_category) VALUES (30, 1, 2);`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert expense;", err)
				return
			}

			// Get accounts and categories
			accounts, err = getAccounts()
			if err != nil {
				return
			}
			categories, err = getCats()
			if err != nil {
				return
			}

			// Check accounts
			if accounts[0].CurrentAmount != 40 || accounts[1].CurrentAmount != 200 {
				t.Errorf("accounts current amount is wrong; expected 40, 200; received: %f, %f", accounts[0].CurrentAmount, accounts[1].CurrentAmount)
				return
			}

			// Check categories
			if categories[0].CurrentAmount != 70 || categories[1].CurrentAmount != 170 {
				t.Errorf("category current amount is wrong; expected 70, 170; received: %f, %f", categories[0].CurrentAmount, categories[1].CurrentAmount)
				return
			}
		},

		// You can't add and Expense that changes the account or category current_amount bellow zero
		func(t *testing.T) {
			// Run init
			err := beforeExpenseTest(t)
			if err != nil {
				t.Error(err)
				return
			}

			// Add expense
			stmt := `INSERT INTO expenses (amount, from_account, from_category) VALUES (150, 1, 2);`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("expected an error; shouldn't be able to reduce account current amount bellow zero", err)
				return
			}

			// Add expense
			stmt = `INSERT INTO expenses (amount, from_account, from_category) VALUES (150, 2, 1);`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("expected an error; shouldn't be able to reduce category current amount bellow zero", err)
				return
			}
		},

		// When you remove an Expense, the amount gets reflected in the related account current_amount and usage_count and in the category current_amount
		func(t *testing.T) {
			// Run init
			err := beforeExpenseTest(t)
			if err != nil {
				t.Error(err)
				return
			}

			// Add expense
			stmt := `INSERT INTO expenses (amount, from_account, from_category) VALUES (10, 1, 1);
					 INSERT INTO expenses (amount, from_account, from_category) VALUES (20, 1, 1);`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert new expenseS", err)
				return
			}

			// Delete expense
			stmt = `DELETE FROM expenses WHERE ID=2`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete expense", err)
				return
			}

			// Get accounts
			getAccounts := func() ([]models.Account, error) {

				// Get accounts
				query := `SELECT current_amount FROM accounts WHERE id=1;`

				// Set variable
				var accounts []models.Account

				// Get rows
				rows, err := db.Query(query)
				if err != nil {
					t.Error("couldn't get account", err)
					return nil, err
				}

				for rows.Next() {
					var account models.Account
					err = rows.Scan(&account.CurrentAmount)
					if err != nil {
						t.Error("couldn't scan account", err)
						return nil, err
					}
					accounts = append(accounts, account)
				}

				return accounts, nil
			}

			accounts, err := getAccounts()
			if err != nil {
				return
			}

			// Check account current amount
			if accounts[0].CurrentAmount != 90 {
				t.Errorf("account current amount is wrong; expected 90; received: %f", accounts[0].CurrentAmount)
				return
			}

			getCats := func() ([]models.Category, error) {
				// Get categories
				query := `SELECT current_amount FROM categories WHERE id=1;`

				// Set variable
				var categories []models.Category

				// Get rows
				rows, err := db.Query(query)
				if err != nil {
					t.Error("couldn't get categories", err)
					return nil, err
				}

				for rows.Next() {
					var category models.Category
					err = rows.Scan(&category.CurrentAmount)
					if err != nil {
						t.Error("couldn't scan category", err)
						return nil, err
					}
					categories = append(categories, category)
				}

				return categories, nil
			}

			categories, err := getCats()
			if err != nil {
				return
			}

			// Check if account current amount is changed
			if categories[0].CurrentAmount != 90 {
				t.Errorf("category current amount is wrong; expected 90; received: %f", categories[0].CurrentAmount)
				return
			}
		},

		// When you update an Expense amount, the amount gets reflected in the related account current_amount and usage_count and in the category current_amount
		func(t *testing.T) {
			// Run init
			err := beforeExpenseTest(t)
			if err != nil {
				t.Error(err)
				return
			}

			// Add expense
			stmt := `INSERT INTO expenses (amount, from_account, from_category) VALUES (10, 1, 1);`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert new expenseS", err)
				return
			}

			// Get accounts
			getAccounts := func() ([]models.Account, error) {

				// Get accounts
				query := `SELECT current_amount, usage_count FROM accounts`

				// Set variable
				var accounts []models.Account

				// Get rows
				rows, err := db.Query(query)
				if err != nil {
					t.Error("couldn't get account", err)
					return nil, err
				}

				for rows.Next() {
					var account models.Account
					err = rows.Scan(&account.CurrentAmount, &account.UsageCount)
					if err != nil {
						t.Error("couldn't scan account", err)
						return nil, err
					}
					accounts = append(accounts, account)
				}

				return accounts, nil
			}

			accounts, err := getAccounts()
			if err != nil {
				return
			}

			// Check account current amount
			if accounts[0].CurrentAmount != 90 {
				t.Errorf("account current amount is wrong; expected 90; received: %f", accounts[0].CurrentAmount)
				return
			}
			if accounts[0].UsageCount != 1 {
				t.Errorf("account usage_count is wrong; expected 1; received: %d", accounts[0].UsageCount)
				return
			}

			// Update expense
			stmt = `UPDATE expenses SET amount = 20 WHERE id=1`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't update expense", err)
				return
			}

			accounts, err = getAccounts()
			if err != nil {
				return
			}

			// Check account current amount
			if accounts[0].CurrentAmount != 80 {
				t.Errorf("account current amount is wrong; expected 80; received: %f", accounts[0].CurrentAmount)
				return
			}
			if accounts[0].UsageCount != 1 {
				t.Errorf("account usage_count is wrong; expected 1; received: %d", accounts[0].UsageCount)
				return
			}

			getCats := func() ([]models.Category, error) {
				// Get categories
				query := `SELECT current_amount FROM categories WHERE id=1;`

				// Set variable
				var categories []models.Category

				// Get rows
				rows, err := db.Query(query)
				if err != nil {
					t.Error("couldn't get categories", err)
					return nil, err
				}

				for rows.Next() {
					var category models.Category
					err = rows.Scan(&category.CurrentAmount)
					if err != nil {
						t.Error("couldn't scan category", err)
						return nil, err
					}
					categories = append(categories, category)
				}

				return categories, nil
			}

			categories, err := getCats()
			if err != nil {
				return
			}

			// Check if account current amount is changed
			if categories[0].CurrentAmount != 80 {
				t.Errorf("category current amount is wrong; expected 80; received: %f", categories[0].CurrentAmount)
				return
			}
		},

		// When you update an Expense account and category, the amount gets reflected in the related accounts current_amount and usage_count and in the categories current_amount
		func(t *testing.T) {
			// Run init
			err := beforeExpenseTest(t)
			if err != nil {
				t.Error(err)
				return
			}

			// Add expense
			stmt := `INSERT INTO expenses (amount, from_account, from_category) VALUES (10, 1, 1);`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert new expense", err)
				return
			}

			// Get accounts
			getAccounts := func() ([]models.Account, error) {

				// Get accounts
				query := `SELECT current_amount, usage_count FROM accounts`

				// Set variable
				var accounts []models.Account

				// Get rows
				rows, err := db.Query(query)
				if err != nil {
					t.Error("couldn't get account", err)
					return nil, err
				}

				for rows.Next() {
					var account models.Account
					err = rows.Scan(&account.CurrentAmount, &account.UsageCount)
					if err != nil {
						t.Error("couldn't scan account", err)
						return nil, err
					}
					accounts = append(accounts, account)
				}

				return accounts, nil
			}

			accounts, err := getAccounts()
			if err != nil {
				return
			}

			// Check account current amount
			if accounts[0].CurrentAmount != 90 || accounts[1].CurrentAmount != 200 {
				t.Errorf("account current amount is wrong; expected 90, 200; received: %f, %f", accounts[0].CurrentAmount, accounts[1].CurrentAmount)
				return
			}
			if accounts[0].UsageCount != 1 || accounts[1].UsageCount != 0 {
				t.Errorf("account usage_count is wrong; expected 1, 0; received: %d, %d", accounts[0].UsageCount, accounts[1].UsageCount)
				return
			}

			// Update expense
			stmt = `UPDATE expenses SET from_account = 2, from_category = 2 WHERE id=1`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't update expense", err)
				return
			}

			accounts, err = getAccounts()
			if err != nil {
				return
			}

			// Check account current amount
			if accounts[0].CurrentAmount != 100 || accounts[1].CurrentAmount != 190 {
				t.Errorf("account current amount is wrong; expected 100, 190; received: %f, %f", accounts[0].CurrentAmount, accounts[1].CurrentAmount)
				return
			}
			if accounts[0].UsageCount != 0 || accounts[1].UsageCount != 1 {
				t.Errorf("account usage_count is wrong; expected 0, 1; received: %d, %d", accounts[0].UsageCount, accounts[1].UsageCount)
				return
			}

			getCats := func() ([]models.Category, error) {
				// Get categories
				query := `SELECT current_amount FROM categories WHERE id=1;`

				// Set variable
				var categories []models.Category

				// Get rows
				rows, err := db.Query(query)
				if err != nil {
					t.Error("couldn't get categories", err)
					return nil, err
				}

				for rows.Next() {
					var category models.Category
					err = rows.Scan(&category.CurrentAmount)
					if err != nil {
						t.Error("couldn't scan category", err)
						return nil, err
					}
					categories = append(categories, category)
				}

				return categories, nil
			}

			categories, err := getCats()
			if err != nil {
				return
			}

			// Check if account current amount is changed
			if categories[0].CurrentAmount != 100 || categories[1].CurrentAmount != 190 {
				t.Errorf("category current amount is wrong; expected 100, 190; received: %f, %f", categories[0].CurrentAmount, categories[1].CurrentAmount)
				return
			}
		},

		/*
		 * Test Data Pipelines
		 *
		 */
		// Test category input period reset procedure
	)
}

func beforeExpenseTest(t *testing.T) error {
	// Add accounts
	stmt := `INSERT INTO accounts (name) VALUES ('test account1');
			 INSERT INTO accounts (name) VALUES ('test account2');`

	// Execute
	_, err := db.Exec(stmt)
	if err != nil {
		t.Error("couldn't insert accounts;", err)
		return err
	}

	// Add categories
	stmt = `INSERT INTO categories (name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category1', 100, 1, 2, 100, 1, 2);
			INSERT INTO categories (name, budget_input, input_interval, input_period, spending_limit, spending_interval, spending_period) VALUES ('test category2', 100, 1, 2, 100, 1, 2);`

	// Execute
	_, err = db.Exec(stmt)
	if err != nil {
		t.Error("couldn't insert categories;", err)
		return err
	}

	// Add money to categories and accounts
	stmt = `UPDATE categories SET initial_amount = 100 WHERE id=1;
			UPDATE categories SET initial_amount = 200 WHERE id=2;
			UPDATE accounts	  SET initial_amount = 100 WHERE id=1;
			UPDATE accounts	  SET initial_amount = 200 WHERE id=2;`

	// Execute
	_, err = db.Exec(stmt)
	if err != nil {
		t.Error("couldn't insert categories;", err)
		return err
	}

	// Add tags
	stmt = `INSERT INTO tags (name) VALUES ('tag 1');
			INSERT INTO tags (name) VALUES ('tag 2');`

	// Execute
	_, err = db.Exec(stmt)
	if err != nil {
		t.Error("couldn't insert categories;", err)
		return err
	}

	return nil
}
