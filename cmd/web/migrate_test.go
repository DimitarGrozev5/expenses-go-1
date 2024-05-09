package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
	cp "github.com/otiai10/copy"
)

var db *sql.DB

// Setup functions
func beforeAll(t *testing.T) {
	if db == nil {
		// Get DB name
		dbName := "test"

		// Remove old db
		os.Remove("test.db")
		os.Remove("test_copy.db")

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

		// Insert user
		stmt := `INSERT INTO user (email, password, db_version) VALUES ('test@test.test', 'adsdasdsadas', 1)`
		_, err = db.Exec(stmt)
		if err != nil {
			t.Error("couldn't add free funds;", err)
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
	os.Remove("test_copy.db")
}

// Setup multiple tests
func setupMultiple(t *testing.T, tests ...func(t *testing.T)) {
	beforeAll(t)
	defer afterAll()

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
		// You can insert an account using an insert procedure
		func(t *testing.T) {
			// Insert account
			stmt := `INSERT INTO procedure_insert_account (name) VALUES ('test account')`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert an account using procedure", err)
				return
			}

			// Get accounts
			accounts, err := getAccounts(t)
			if err != nil {
				return
			}

			// Check if order is ok
			if len(accounts) > 1 {
				t.Errorf("too many accounts: %d, expected 1", len(accounts))
				return
			}
			if accounts[0].Name != "test account" {
				t.Errorf("wrong account name; expected 'test account'; received %s", accounts[0].Name)
				return
			}
			if accounts[0].CurrentAmount != 0 {
				t.Errorf("wrong account current amount; expected 0; received %.2f", accounts[0].CurrentAmount)
				return
			}
			if accounts[0].UsageCount != 0 {
				t.Errorf("wrong account usage count; expected 0; received %d", accounts[0].UsageCount)
				return
			}
			if accounts[0].TableOrder != 1 {
				t.Errorf("wrong account table order; expected 1; received %d", accounts[0].TableOrder)
				return
			}

			// Insert account
			stmt = `INSERT INTO procedure_insert_account (name) VALUES ('test account 1')`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert an account using procedure", err)
				return
			}

			// Get accounts
			accounts, err = getAccounts(t)
			if err != nil {
				return
			}

			// Check if order is ok
			if len(accounts) > 2 {
				t.Errorf("too many accounts: %d, expected 2", len(accounts))
				return
			}
			if accounts[1].Name != "test account 1" {
				t.Errorf("wrong account name; expected 'test account 1'; received %s", accounts[1].Name)
				return
			}
			if accounts[1].CurrentAmount != 0 {
				t.Errorf("wrong account current amount; expected 0; received %.2f", accounts[1].CurrentAmount)
				return
			}
			if accounts[1].UsageCount != 0 {
				t.Errorf("wrong account usage count; expected 0; received %d", accounts[1].UsageCount)
				return
			}
			if accounts[1].TableOrder != 2 {
				t.Errorf("wrong account table order; expected 2; received %d", accounts[1].TableOrder)
				return
			}
		},

		// You can rename account using procedure
		func(t *testing.T) {
			// Insert account
			stmt := `INSERT INTO procedure_insert_account (name) VALUES ('test account');
					 INSERT INTO procedure_insert_account (name) VALUES ('test account2')`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert an account using procedure", err)
				return
			}

			// Get accounts
			accounts, err := getAccounts(t)
			if err != nil {
				return
			}

			// Check if order is ok
			if accounts[0].Name != "test account" {
				t.Errorf("wrong account name; expected 'test account'; received %s", accounts[0].Name)
				return
			}

			// Rename account
			stmt = `UPDATE procedure_account_update_name SET name='test name 2' WHERE id=1`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert an account using procedure", err)
				return
			}

			// Get accounts
			accounts, err = getAccounts(t)
			if err != nil {
				return
			}

			// Check name
			if accounts[0].Name != "test name 2" || accounts[1].Name != "test account2" {
				t.Errorf("wrong account name; expected 'test name 2', 'test account2'; received '%s', '%s'", accounts[0].Name, accounts[1].Name)
				return
			}
		},

		// You can move account up and down
		func(t *testing.T) {
			// Insert accounts
			stmt := `	INSERT INTO procedure_insert_account (name) VALUES ('acc 1');
						INSERT INTO procedure_insert_account (name) VALUES ('acc 2');
						INSERT INTO procedure_insert_account (name) VALUES ('acc 3')`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert a accounts using procedure", err)
				return
			}

			testOrder := func(a, b, c int) error {

				// Get accounts
				accounts, err := getAccounts(t)
				if err != nil {
					return err
				}

				// Check if order is ok
				if !(accounts[0].TableOrder == a && accounts[1].TableOrder == b && accounts[2].TableOrder == c) {
					t.Errorf(
						"error with accounts order; expected %d, %d, %d; received %d, %d, %d",
						a,
						b,
						c,
						accounts[0].TableOrder,
						accounts[1].TableOrder,
						accounts[2].TableOrder,
					)
					return errors.New("error with acount order")
				}

				return nil
			}

			// Test order
			err = testOrder(1, 2, 3)
			if err != nil {
				return
			}

			////////////////////////////////////////////////////////////// Move first account up
			stmt = `UPDATE procedure_change_accounts_order SET table_order=4 WHERE id=3;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("should not be able to move first account up; expected error; received nil")
				return
			}
			if !strings.HasPrefix(err.Error(), "cant move first account up") {
				t.Errorf("wrong error; expected 'cant move first account up'; received %s", err)
				return
			}

			// Test order
			err = testOrder(1, 2, 3)
			if err != nil {
				return
			}

			////////////////////////////////////////////////////////////// Move last account down
			stmt = `UPDATE procedure_change_accounts_order SET table_order=0 WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("should not be able to move last account down; expected error; received nil")
				return
			}
			if !strings.HasPrefix(err.Error(), "cant move last account down") {
				t.Errorf("wrong error; expected 'cant move last account down'; received %s", err.Error())
				return
			}

			// Test order
			err = testOrder(1, 2, 3)
			if err != nil {
				return
			}

			////////////////////////////////////////////////////////////// Move first account to the middle
			stmt = `UPDATE procedure_change_accounts_order SET table_order=2 WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't change account table order", err)
				return
			}

			// Test order
			err = testOrder(2, 1, 3)
			if err != nil {
				return
			}

			////////////////////////////////////////////////////////////// Move first account back to the start
			stmt = `UPDATE procedure_change_accounts_order SET table_order=1 WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't change account table order", err)
				return
			}

			// Test order
			err = testOrder(1, 2, 3)
			if err != nil {
				return
			}
		},

		// You can delete an unused account and the table_order will update on all others
		func(t *testing.T) {
			// Insert accounts
			stmt := `	INSERT INTO procedure_insert_account (name) VALUES ('acc 1');
						INSERT INTO procedure_insert_account (name) VALUES ('acc 2');
						INSERT INTO procedure_insert_account (name) VALUES ('acc 3');
						INSERT INTO procedure_insert_account (name) VALUES ('acc 4');`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert a accounts using procedure", err)
				return
			}

			testOrder := func(acc ...[]int) error {

				// Get accounts
				accounts, err := getAccounts(t)
				if err != nil {
					return err
				}

				for i, a := range acc {

					// Check if order is ok
					if accounts[i].TableOrder != a[1] || accounts[i].ID != a[0] {
						t.Errorf(
							"error with accounts order; expected %d, %d; received %d, %d",
							a[0],
							a[1],
							accounts[i].ID,
							accounts[i].TableOrder,
						)
						return errors.New("error with acount order")
					}
				}

				return nil
			}

			err = testOrder([]int{1, 1}, []int{2, 2}, []int{3, 3}, []int{4, 4})
			if err != nil {
				return
			}

			// Delete second account
			stmt = `DELETE FROM procedure_remove_account WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete accounts using procedure", err)
				return
			}

			// Test order
			err = testOrder([]int{1, 1}, []int{3, 2}, []int{4, 3})
			if err != nil {
				return
			}

			// Delete third account
			stmt = `DELETE FROM procedure_remove_account WHERE id=3;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete accounts using procedure", err)
				return
			}

			// Test order
			err = testOrder([]int{1, 1}, []int{4, 2})
			if err != nil {
				return
			}
		},

		// You can't delete an account that is being used. Account is being used when it's referenced by excpenses or when it has funds
		func(t *testing.T) {
			// Insert accounts
			stmt := `	INSERT INTO procedure_insert_account (name) VALUES ('acc 1');
						INSERT INTO procedure_insert_account (name) VALUES ('acc 2');`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert a accounts using procedure", err)
				return
			}

			// Artificialy update account current amount
			stmt = `UPDATE accounts SET current_amount=100 WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't update account", err)
				return
			}

			// Try to delete first account
			stmt = `DELETE FROM procedure_remove_account WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("you shouldn't be able to delete an account that is used", err)
				return
			}
			if !strings.HasPrefix(err.Error(), "cant delete an account that is used") {
				t.Errorf("wrong error received; expected 'cant delete an account that is used'; received %s", err)
				return
			}

			// Try to delete second account
			stmt = `DELETE FROM procedure_remove_account WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("failed to delete an account that is not being used", err)
				return
			}
		},

		/*
		 * Categories tests
		 *
		 */
		// You can insert a category using an insert procedure
		func(t *testing.T) {
			// Insert category
			stmt := `INSERT INTO procedure_new_category (
						name,
						budget_input,
						input_interval,
						input_period,
						spending_limit
					) VALUES (
						'test category',
						100,
						1,
						1,
						100
					)`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert a category using procedure", err)
				return
			}

			// Get categories
			categories, err := getCategories(t)
			if err != nil {
				return
			}

			// Check if all is ok
			if len(categories) != 1 {
				t.Errorf("too many categories: %d, expected 1", len(categories))
				return
			}
			if categories[0].Name != "test category" {
				t.Errorf("wrong category name; expected 'test category'; received %s", categories[0].Name)
				return
			}

			// Insert account
			stmt = `INSERT INTO procedure_new_category (
						name,
						budget_input,
						input_interval,
						input_period,
						spending_limit
					) VALUES (
						'test category 1',
						100,
						1,
						1,
						100
					)`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert a category using procedure", err)
				return
			}

			// Get category
			categories, err = getCategories(t)
			if err != nil {
				return
			}

			// TODO: Add tests for all fields of the categories table
			// Check if all is ok
			if len(categories) != 2 {
				t.Errorf("too many categories: %d, expected 2", len(categories))
				return
			}
			if categories[1].Name != "test category 1" {
				t.Errorf("wrong category name; expected 'test category 1'; received %s", categories[0].Name)
				return
			}
		},

		// You can rename category using procedure
		func(t *testing.T) {
			// Insert category
			stmt := `INSERT INTO procedure_new_category (
						name,
						budget_input,
						input_interval,
						input_period,
						spending_limit
					) VALUES (
						'test category',
						100,
						1,
						1,
						100
					)`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert a category using procedure", err)
				return
			}

			// Get categories
			categories, err := getCategories(t)
			if err != nil {
				return
			}

			// Check if name is the same
			if categories[0].Name != "test category" {
				t.Errorf("wrong category name; expected 'test category'; received %s", categories[0].Name)
				return
			}

			// Rename category
			stmt = `UPDATE procedure_category_name SET name='test name 2' WHERE id=1`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't update category name using procedure", err)
				return
			}

			// Get categories
			categories, err = getCategories(t)
			if err != nil {
				return
			}

			// Check if name is the same
			if categories[0].Name != "test name 2" {
				t.Errorf("wrong category name; expected 'test name 2'; received %s", categories[0].Name)
				return
			}
		},

		// You can move category up and down
		func(t *testing.T) {
			// Insert category
			stmt := `INSERT INTO procedure_new_category (
				name,
				budget_input,
				input_interval,
				input_period,
				spending_limit
			) VALUES (
				'test category1',
				100,
				1,
				1,
				100
			);
			INSERT INTO procedure_new_category (
				name,
				budget_input,
				input_interval,
				input_period,
				spending_limit
			) VALUES (
				'test category2',
				100,
				1,
				1,
				100
			);
			INSERT INTO procedure_new_category (
				name,
				budget_input,
				input_interval,
				input_period,
				spending_limit
			) VALUES (
				'test category3',
				100,
				1,
				1,
				100
			);`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert categories using procedure", err)
				return
			}

			testOrder := func(a, b, c int) error {

				// Get categories
				categories, err := getCategories(t)
				if err != nil {
					return err
				}

				// Check if order is ok
				if !(categories[0].TableOrder == a && categories[1].TableOrder == b && categories[2].TableOrder == c) {
					t.Errorf(
						"error with categories order; expected %d, %d, %d; received %d, %d, %d",
						a,
						b,
						c,
						categories[0].TableOrder,
						categories[1].TableOrder,
						categories[2].TableOrder,
					)
					return errors.New("error with category order")
				}

				return nil
			}

			// Test order
			err = testOrder(1, 2, 3)
			if err != nil {
				return
			}

			////////////////////////////////////////////////////////////// Move first category up
			stmt = `UPDATE procedure_change_categories_order SET table_order=4 WHERE id=3;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("should not be able to move first category up; expected error; received nil")
				return
			}
			if !strings.HasPrefix(err.Error(), "cant move first category up") {
				t.Errorf("wrong error; expected 'cant move first category up'; received %s", err)
				return
			}

			// Test order
			err = testOrder(1, 2, 3)
			if err != nil {
				return
			}

			////////////////////////////////////////////////////////////// Move last category down
			stmt = `UPDATE procedure_change_categories_order SET table_order=0 WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("should not be able to move last category down; expected error; received nil")
				return
			}
			if !strings.HasPrefix(err.Error(), "cant move last category down") {
				t.Errorf("wrong error; expected 'cant move last category down'; received %s", err.Error())
				return
			}

			// Test order
			err = testOrder(1, 2, 3)
			if err != nil {
				return
			}

			////////////////////////////////////////////////////////////// Move first category to the middle
			stmt = `UPDATE procedure_change_categories_order SET table_order=2 WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't change category table order", err)
				return
			}

			// Test order
			err = testOrder(2, 1, 3)
			if err != nil {
				return
			}

			////////////////////////////////////////////////////////////// Move first category back to the start
			stmt = `UPDATE procedure_change_categories_order SET table_order=1 WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't change category table order", err)
				return
			}

			// Test order
			err = testOrder(1, 2, 3)
			if err != nil {
				return
			}
		},

		// You can delete an unused category and the table_order will update on all others
		func(t *testing.T) {
			// Insert category
			stmt := `INSERT INTO procedure_new_category (
				name,
				budget_input,
				input_interval,
				input_period,
				spending_limit
			) VALUES (
				'test category1',
				100,
				1,
				1,
				100
			);
			INSERT INTO procedure_new_category (
				name,
				budget_input,
				input_interval,
				input_period,
				spending_limit
			) VALUES (
				'test category2',
				100,
				1,
				1,
				100
			);
			INSERT INTO procedure_new_category (
				name,
				budget_input,
				input_interval,
				input_period,
				spending_limit
			) VALUES (
				'test category3',
				100,
				1,
				1,
				100
			);
			INSERT INTO procedure_new_category (
				name,
				budget_input,
				input_interval,
				input_period,
				spending_limit
			) VALUES (
				'test category4',
				100,
				1,
				1,
				100
			);`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert categories using procedure", err)
				return
			}

			testOrder := func(acc ...[]int) error {

				// Get categories
				categories, err := getCategories(t)
				if err != nil {
					t.Error("couldn't get categories")
					return err
				}

				for i, a := range acc {

					// Check if order is ok
					if categories[i].TableOrder != a[1] || categories[i].ID != a[0] {
						t.Errorf(
							"error with categories order; expected %d, %d; received %d, %d",
							a[0],
							a[1],
							categories[i].ID,
							categories[i].TableOrder,
						)
						return errors.New("error with category order")
					}
				}

				return nil
			}

			err = testOrder([]int{1, 1}, []int{2, 2}, []int{3, 3}, []int{4, 4})
			if err != nil {
				return
			}

			// Delete second category
			stmt = `DELETE FROM procedure_remove_category WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete categories using procedure", err)
				return
			}

			// Test order
			err = testOrder([]int{1, 1}, []int{3, 2}, []int{4, 3})
			if err != nil {
				return
			}

			// Delete third category
			stmt = `DELETE FROM procedure_remove_category WHERE id=3;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete categories using procedure", err)
				return
			}

			// Test order
			err = testOrder([]int{1, 1}, []int{4, 2})
			if err != nil {
				return
			}
		},

		// You can't delte a cateogory that is being used. Category is being used when it's referenced by expenses or when it has funds
		func(t *testing.T) {
			// Insert category
			stmt := `INSERT INTO procedure_new_category (
				name,
				budget_input,
				input_interval,
				input_period,
				spending_limit
			) VALUES (
				'test category1',
				100,
				1,
				1,
				100
			);
			INSERT INTO procedure_new_category (
				name,
				budget_input,
				input_interval,
				input_period,
				spending_limit
			) VALUES (
				'test category2',
				100,
				1,
				1,
				100
			);`

			// Execute
			_, err := db.Exec(stmt)
			if err != nil {
				t.Error("couldn't insert categories using procedure", err)
				return
			}

			// Artificialy update category current amount
			stmt = `UPDATE categories SET initial_amount=100 WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't update category", err)
				return
			}

			// Try to delete first category
			stmt = `DELETE FROM procedure_remove_category WHERE id=1;`

			// Execute
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("you shouldn't be able to delete a category that is used", err)
				return
			}
			if !strings.HasPrefix(err.Error(), "cant delete a category that is used") {
				t.Errorf("wrong error received; expected 'cant delete a category that is used'; received %s", err)
				return
			}

			// Try to delete second category
			stmt = `DELETE FROM procedure_remove_category WHERE id=2;`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("failed to delete an category that is not being used", err)
				return
			}
		},

		/*
		 * Test Expense interactions
		 *
		 */
		// You can insert expenses, tags and expense tags, using a procedure
		func(t *testing.T) {
			// Run init
			err := beforeExpenseTest(t)
			if err != nil {
				t.Error(err)
				return
			}

			// Insert an expense
			stmt := `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (10, $1, 1, 1);`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't insert expense using procedure", err)
				return
			}

			// Get expenses
			expenses, err := getExpenses(t)
			if err != nil {
				return
			}

			if len(expenses) != 1 {
				t.Errorf("too many expenses; expected 1; received %d", len(expenses))
				return
			}

			// Insert new tag
			stmt = `INSERT INTO procedure_insert_tag (name) VALUES ('test new tag');`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't insert tag using procedure", err)
				return
			}

			// Get tags
			tags, err := getTags(t)
			if err != nil {
				return
			}

			// Check tags length
			if len(tags) != 3 {
				t.Errorf("too many tags; expected 3; received %d", len(tags))
				return
			}

			// Insert same tag again (should not throw an error) along with new tags
			stmt = `INSERT INTO procedure_insert_tag (name) VALUES ('test new tag'), ('test new 1'), ('test new 2');`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't insert existing tag using procedure", err)
				return
			}

			// Get tags
			tags, err = getTags(t)
			if err != nil {
				return
			}

			// Check tags length
			if len(tags) != 5 {
				t.Errorf("too many tags; expected 5; received %d", len(tags))
				return
			}

			// Insert expense_tag
			stmt = `INSERT iNTO procedure_link_tag_to_expense (expense_id, tag_id) VALUES (1, 1)`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't insert existing tag using procedure", err)
				return
			}

			// Get relations
			rel, err := getExpenseTags(t)
			if err != nil {
				return
			}

			// Check length
			if len(rel) != 1 {
				t.Errorf("too many expense-tag relations; expected 1; received %d", len(rel))
				return
			}
		},

		// When you add an Expense with tag relations it changes the related tags usage_count value
		// You can't delete tags, that are used
		func(t *testing.T) {
			// Run init
			err := beforeExpenseTest(t)
			if err != nil {
				t.Error(err)
				return
			}

			// Get initial tags
			initialTags, err := getTags(t)
			if err != nil {
				return
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
			stmt := `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (10, $1, 1, 1);
					 INSERT iNTO procedure_link_tag_to_expense (expense_id, tag_id) VALUES (1, 1)`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't insert expense and expense_tags;", err)
				return
			}

			// Get tags
			tags1, err := getTags(t)
			if err != nil {
				return
			}

			// Check tags initial usage_count
			if len(tags1) > 2 {
				t.Errorf("too many tags: %d; expected 2", len(tags1))
				return
			}
			if tags1[0].UsageCount != 1 || tags1[1].UsageCount != 0 {
				t.Errorf("unexpected tags usage_count; recevied: %d, %d; expected 1, 0", tags1[0].UsageCount, tags1[1].UsageCount)
				return
			}

			// Add expense with relation to tag 1 and 2
			stmt = `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (10, $1, 1, 1);
					INSERT iNTO procedure_link_tag_to_expense (expense_id, tag_id) VALUES (2, 1);
					INSERT iNTO procedure_link_tag_to_expense (expense_id, tag_id) VALUES (2, 2);`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't insert expense and expense_tags;", err)
				return
			}

			// Get tags
			tags2, err := getTags(t)
			if err != nil {
				return
			}

			// Check tags initial usage_count
			if len(tags2) > 2 {
				t.Errorf("too many tags: %d; expected 2", len(tags2))
				return
			}
			if tags2[0].UsageCount != 2 || tags2[1].UsageCount != 1 {
				t.Errorf("unexpected tags usage_count; recevied: %d, %d; expected 2, 1", tags2[0].UsageCount, tags2[1].UsageCount)
				return
			}

			// Delete from expense tags
			stmt = `DELETE FROM procedure_unlink_tag_from_expense WHERE id=3`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete expense tags;", err)
				return
			}

			// Get tags
			tags3, err := getTags(t)
			if err != nil {
				return
			}

			// Check tags initial usage_count
			if len(tags3) > 2 {
				t.Errorf("too many tags: %d; expected 2", len(tags3))
				return
			}
			if tags3[0].UsageCount != 2 || tags3[1].UsageCount != 0 {
				t.Errorf("unexpected tags usage_count; recevied: %d, %d; expected 2, 0", tags3[0].UsageCount, tags3[1].UsageCount)
				return
			}

			// Remove expense
			stmt = `DELETE FROM procedure_remove_expense WHERE id=1`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete expense;", err)
				return
			}

			// Get tags
			tags4, err := getTags(t)
			if err != nil {
				return
			}

			// Check tags usage_count
			if len(tags4) > 2 {
				t.Errorf("too many tags: %d; expected 2", len(tags4))
				return
			}
			if tags4[0].UsageCount != 1 || tags4[1].UsageCount != 0 {
				t.Errorf("unexpected tags usage_count; recevied: %d, %d; expected 1, 0", tags4[0].UsageCount, tags4[1].UsageCount)
				return
			}

			// Delete second tag
			stmt = `DELETE FROM procedure_remove_tag WHERE id=2`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete tag;", err)
				return
			}

			// Get tags
			tags5, err := getTags(t)
			if err != nil {
				return
			}

			// Check tags usage_count
			if len(tags5) > 1 {
				t.Errorf("too many tags: %d; expected 1", len(tags5))
				return
			}

			// Delete first tag
			stmt = `DELETE FROM procedure_remove_tag WHERE id=1`

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
			stmt := `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (10, $1, 1, 1);`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't insert expense;", err)
				return
			}

			accounts, err := getAccounts(t)
			if err != nil {
				return
			}

			// Check if account current amount is changed
			if accounts[0].CurrentAmount != 90 {
				t.Errorf("account current amount is wrong; expected 90; received: %f", accounts[0].CurrentAmount)
				return
			}

			categories, err := getCategories(t)
			if err != nil {
				return
			}

			// Check if account current amount is changed
			if categories[0].CurrentAmount != 90 {
				t.Errorf("category current amount is wrong; expected 90; received: %f", categories[0].CurrentAmount)
				return
			}

			// Add expenses
			stmt = `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (20, $1, 1, 1);
					INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (30, $2, 1, 2);`

			// Execute
			_, err = db.Exec(stmt, time.Now(), time.Now())
			if err != nil {
				t.Error("couldn't insert expenses;", err)
				return
			}

			// Get accounts and categories
			accounts, err = getAccounts(t)
			if err != nil {
				return
			}
			categories, err = getCategories(t)
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
			stmt := `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (150, $1, 1, 2);`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err == nil {
				t.Error("expected an error; shouldn't be able to reduce account current amount bellow zero", err)
				return
			}

			// Add expense
			stmt = `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (150, $1, 2, 1);`

			// Execute
			_, err = db.Exec(stmt, time.Now())
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
			stmt := `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (10, $1, 1, 1);
					 INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (20, $2, 1, 1);`

			// Execute
			_, err = db.Exec(stmt, time.Now(), time.Now())
			if err != nil {
				t.Error("couldn't insert new expenses", err)
				return
			}

			// Delete expense
			stmt = `DELETE FROM procedure_remove_expense WHERE ID=2`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't delete expense", err)
				return
			}

			accounts, err := getAccounts(t)
			if err != nil {
				return
			}

			// Check account current amount
			if accounts[0].CurrentAmount != 90 {
				t.Errorf("account current amount is wrong; expected 90; received: %f", accounts[0].CurrentAmount)
				return
			}

			categories, err := getCategories(t)
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
			stmt := `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (10, $1, 1, 1);`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't insert new expense", err)
				return
			}

			accounts, err := getAccounts(t)
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
			stmt = `UPDATE procedure_update_expense SET amount = 20 WHERE id=1`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't update expense", err)
				return
			}

			accounts, err = getAccounts(t)
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

			categories, err := getCategories(t)
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
			stmt := `INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (10, $1, 1, 1);`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't insert new expense", err)
				return
			}

			accounts, err := getAccounts(t)
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
			stmt = `UPDATE procedure_update_expense SET from_account = 2, from_category = 2 WHERE id=1`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't update expense", err)
				return
			}

			accounts, err = getAccounts(t)
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

			categories, err := getCategories(t)
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
		// Test free funds procedure
		// Free funds are set in the user table and the amount is reflected in a target account
		// If the db is used only through procedures then the accounts total amount should never go bellow the free funds
		func(t *testing.T) {
			// Insert accounts and categories
			err := beforeExpenseTest(t)
			if err != nil {
				t.Error(err)
				return
			}

			// Add free funds
			stmt := `INSERT INTO procedure_add_free_funds (amount, to_account) VALUES (100, 1)`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't add free funds;", err)
				return
			}

			// Get accounts
			accounts, err := getAccounts(t)
			if err != nil {
				return
			}

			// Check amount in account 1
			if accounts[0].CurrentAmount != 200 {
				t.Errorf("account current amount is wrong; expected 200; received: %f", accounts[0].CurrentAmount)
				return
			}

			// Get user
			user, err := getUser(t)
			if err != nil {
				return
			}

			// Check amount in account 1
			if user.FreeFunds != 100 {
				t.Errorf("user free funds amount is wrong; expected 100; received: %f", user.FreeFunds)
				return
			}

			// Try to drop free funds bellow zero
			stmt = `INSERT INTO procedure_add_free_funds (amount, to_account) VALUES (-150, 1)`

			// Execute
			_, err = db.Exec(stmt, time.Now())
			if err == nil {
				t.Error("expected error; should'n be able to drop free funds bellow zero")
				return
			}
		},

		// Test adding funds to a category and reseting the input period
		func(t *testing.T) {
			// Seed account and category
			err := seedAccAndCat(t)
			if err != nil {
				return
			}

			// Add money to free funds
			stmt := `INSERT INTO procedure_add_free_funds (amount, to_account) VALUES (100, 1)`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't add free funds;", err)
				return
			}

			// Add too much money to category
			stmt = `INSERT INTO procedure_fund_category_and_reset_period (amount, category) VALUES (200, 1)`
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("shouldn't be able to add to category more money than in free funds")
				return
			}

			//////////////////////////////////////////////////////////////////////////////////////////////// Add normal amount to category
			stmt = `INSERT INTO procedure_fund_category_and_reset_period (amount, category) VALUES (60, 1)`
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't add money to category", err)
				return
			}

			// Get data
			user, err := getUser(t)
			if err != nil {
				return
			}
			accounts, err := getAccounts(t)
			if err != nil {
				return
			}
			categories, err := getCategories(t)
			if err != nil {
				return
			}

			// Free funds have to be reduced
			if user.FreeFunds != 40 {
				t.Errorf("free funds are wrong; expected 40; received %f", user.FreeFunds)
				return
			}

			// Account should have the full sum
			if accounts[0].CurrentAmount != 100 {
				t.Errorf("account current amount is wrong; expected 100; received %f", accounts[0].CurrentAmount)
				return
			}

			// Category should have the added amount
			if categories[0].InitialAmount != 60 {
				t.Errorf("category initial amount is wrong; expected 60; received %f", categories[0].InitialAmount)
				return
			}
			if categories[0].SpendingLeft != 80 {
				t.Errorf("category spending left is wrong; expected 80; received %f", categories[0].SpendingLeft)
				return
			}

			//////////////////////////////////////////////////////////////////////////////////////////////// Add expenses and reset
			// Add expense
			stmt = `
					INSERT INTO procedure_insert_tag (name) VALUES ('tag1'), ('tag2');
					INSERT INTO procedure_new_expense (amount, date, from_account, from_category) VALUES (10, $1, 1, 1);
					INSERT INTO procedure_link_tag_to_expense (expense_id, tag_id) VALUES (1, 1), (1, 2);`
			_, err = db.Exec(stmt, time.Now())
			if err != nil {
				t.Error("couldn't add expense", err)
				return
			}

			// Get category
			categories, err = getCategories(t)
			if err != nil {
				return
			}

			// Category should have the updated amount
			if categories[0].InitialAmount != 60 {
				t.Errorf("category initial amount is wrong; expected 60; received %f", categories[0].InitialAmount)
				return
			}
			if categories[0].CurrentAmount != 50 {
				t.Errorf("category current amount is wrong; expected 50; received %f", categories[0].CurrentAmount)
				return
			}
			if categories[0].SpendingLeft != 70 {
				t.Errorf("category spending left is wrong; expected 70; received %f", categories[0].SpendingLeft)
				return
			}

			// Reset category again
			stmt = `INSERT INTO procedure_fund_category_and_reset_period (amount, category) VALUES (20, 1)`
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't add money to category", err)
				return
			}

			// Get data
			user, err = getUser(t)
			if err != nil {
				return
			}
			accounts, err = getAccounts(t)
			if err != nil {
				return
			}
			categories, err = getCategories(t)
			if err != nil {
				return
			}

			// Free funds have to be reduced
			if user.FreeFunds != 20 {
				t.Errorf("free funds are wrong; expected 20; received %f", user.FreeFunds)
				return
			}

			// Account should have the full sum
			if accounts[0].CurrentAmount != 90 {
				t.Errorf("account current amount is wrong; expected 90; received %f", accounts[0].CurrentAmount)
				return
			}

			// Category should have the added amount
			if categories[0].InitialAmount != 70 {
				t.Errorf("category current amount is wrong; expected 70; received %f", categories[0].InitialAmount)
				return
			}
			if categories[0].CurrentAmount != 70 {
				t.Errorf("category current amount is wrong; expected 70; received %f", categories[0].CurrentAmount)
				return
			}
			if categories[0].SpendingLeft != 80 {
				t.Errorf("category spending left is wrong; expected 80; received %f", categories[0].SpendingLeft)
				return
			}

			// Get current expenses
			var currentExpensesCount int
			query := `SELECT COUNT(*) FROM view_current_expenses;`
			row := db.QueryRow(query)
			err = row.Scan(&currentExpensesCount)
			if err != nil {
				t.Error("couldn't get expenses count", err)
				return
			}

			// Check count
			if currentExpensesCount != 0 {
				t.Errorf("too many current expenses; expected 0; received %d", currentExpensesCount)
				return
			}

			//////////////////////////////////////////////////////////////////////////////////////////////// Try to delete archived expense
			stmt = `DELETE FROM procedure_remove_expense WHERE id=1;`
			_, err = db.Exec(stmt)
			if err == nil {
				t.Error("shouldn't be able to delete archived expense")
				return
			}
			if !strings.HasPrefix(err.Error(), "cant delete archived expense") {
				t.Errorf("wrong error; expected: 'cant delete archived expense', received: '%s'", err.Error())
				return
			}

			//////////////////////////////////////////////////////////////////////////////////////////////// Look at archived periods
			var numberOfPeriods int
			query = `SELECT COUNT(*) FROM archived_periods;`
			row = db.QueryRow(query)
			err = row.Scan(&numberOfPeriods)
			if err != nil {
				t.Error("couldn't get expenses count", err)
				return
			}

			// Check count
			if numberOfPeriods != 2 {
				t.Errorf("wrong number of periods; expected 2; received %d", numberOfPeriods)
			}
		},

		// Money can be added and removed from category without reseting the input period
		func(t *testing.T) {
			// Seed account and category
			err := seedAccAndCat(t)
			if err != nil {
				return
			}

			// Add money to free funds
			stmt := `INSERT INTO procedure_add_free_funds (amount, to_account) VALUES (200, 1)`

			// Execute
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't add free funds;", err)
				return
			}

			//////////////////////////////////////////////////////////////////////////////////////////////// Add normal amount to categories
			stmt = `INSERT INTO procedure_fund_category_and_reset_period (amount, category) VALUES (50, 1);
					INSERT INTO procedure_fund_category_and_reset_period (amount, category) VALUES (50, 2)`
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't add money to category", err)
				return
			}

			// Get data
			user, err := getUser(t)
			if err != nil {
				return
			}
			accounts, err := getAccounts(t)
			if err != nil {
				return
			}
			categories, err := getCategories(t)
			if err != nil {
				return
			}

			// Free funds have to be reduced
			if user.FreeFunds != 100 {
				t.Errorf("free funds are wrong; expected 100; received %f", user.FreeFunds)
				return
			}

			// Account should have the full sum
			if accounts[0].CurrentAmount != 200 {
				t.Errorf("account current amount is wrong; expected 200; received %f", accounts[0].CurrentAmount)
				return
			}

			// Category should have the added amount
			if categories[0].InitialAmount != 50 || categories[1].InitialAmount != 50 {
				t.Errorf("category initial amount is wrong; expected 50, 50; received %f, %f", categories[0].InitialAmount, categories[1].InitialAmount)
				return
			}
			if categories[0].SpendingLeft != 80 || categories[1].SpendingLeft != 80 {
				t.Errorf("category spending left is wrong; expected 80, 80; received %f, %f", categories[0].SpendingLeft, categories[1].SpendingLeft)
				return
			}

			//////////////////////////////////////////////////////////////////////////////////////////////// Take money from first category
			stmt = `INSERT INTO procedure_fund_category (amount, category) VALUES (-10, 1);`
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't add money to category", err)
				return
			}

			// Get data
			user, err = getUser(t)
			if err != nil {
				return
			}
			accounts, err = getAccounts(t)
			if err != nil {
				return
			}
			categories, err = getCategories(t)
			if err != nil {
				return
			}

			// Free funds have to be reduced
			if user.FreeFunds != 110 {
				t.Errorf("free funds are wrong; expected 110; received %f", user.FreeFunds)
				return
			}

			// Account should have the full sum
			if accounts[0].CurrentAmount != 200 {
				t.Errorf("account current amount is wrong; expected 200; received %f", accounts[0].CurrentAmount)
				return
			}

			// Category should have the added amount
			if categories[0].InitialAmount != 40 || categories[1].InitialAmount != 50 {
				t.Errorf("category initial amount is wrong; expected 40, 50; received %f, %f", categories[0].InitialAmount, categories[1].InitialAmount)
				return
			}
			if categories[0].SpendingLeft != 80 || categories[1].SpendingLeft != 80 {
				t.Errorf("category spending left is wrong; expected 80, 80; received %f, %f", categories[0].SpendingLeft, categories[1].SpendingLeft)
				return
			}

			//////////////////////////////////////////////////////////////////////////////////////////////// Add money to second category
			stmt = `INSERT INTO procedure_fund_category (amount, category) VALUES (10, 2);`
			_, err = db.Exec(stmt)
			if err != nil {
				t.Error("couldn't add money to category", err)
				return
			}

			// Get data
			user, err = getUser(t)
			if err != nil {
				return
			}
			accounts, err = getAccounts(t)
			if err != nil {
				return
			}
			categories, err = getCategories(t)
			if err != nil {
				return
			}

			// Free funds have to be reduced
			if user.FreeFunds != 100 {
				t.Errorf("free funds are wrong; expected 100; received %f", user.FreeFunds)
				return
			}

			// Account should have the full sum
			if accounts[0].CurrentAmount != 200 {
				t.Errorf("account current amount is wrong; expected 200; received %f", accounts[0].CurrentAmount)
				return
			}

			// Category should have the added amount
			if categories[0].InitialAmount != 40 || categories[1].InitialAmount != 60 {
				t.Errorf("category initial amount is wrong; expected 40, 60; received %f, %f", categories[0].InitialAmount, categories[1].InitialAmount)
				return
			}
			if categories[0].SpendingLeft != 80 || categories[1].SpendingLeft != 80 {
				t.Errorf("category spending left is wrong; expected 80, 80; received %f, %f", categories[0].SpendingLeft, categories[1].SpendingLeft)
				return
			}
		},
	)
}

func seedAccAndCat(t *testing.T) error {
	stmt := `	INSERT INTO accounts (name, table_order) VALUES ('test account1', 1);
				INSERT INTO accounts (name, table_order) VALUES ('test account2', 1);
				INSERT INTO categories (name, budget_input, input_interval, input_period, spending_limit, spending_left, table_order, initial_amount, current_amount) VALUES ('test category1', 100, 1, 2, 80, 80, 1, 0, 0);
				INSERT INTO categories (name, budget_input, input_interval, input_period, spending_limit, spending_left, table_order, initial_amount, current_amount) VALUES ('test category2', 100, 1, 2, 80, 80, 1, 0, 0);`

	// Execute
	_, err := db.Exec(stmt)
	if err != nil {
		t.Error("couldn't insert accounts;", err)
		return err
	}

	return nil
}

func beforeExpenseTest(t *testing.T) error {
	// Add accounts
	stmt := `INSERT INTO accounts (name, current_amount, table_order) VALUES ('test account1', 100, 1);
			 INSERT INTO accounts (name, current_amount, table_order) VALUES ('test account2', 200, 2);`

	// Execute
	_, err := db.Exec(stmt)
	if err != nil {
		t.Error("couldn't insert accounts;", err)
		return err
	}

	// Add categories
	stmt = `INSERT INTO categories (name, budget_input, input_interval, input_period, spending_limit, spending_left, table_order, initial_amount, current_amount) VALUES ('test category1', 100, 1, 2, 100, 100, 1, 100, 100);
			INSERT INTO categories (name, budget_input, input_interval, input_period, spending_limit, spending_left, table_order, initial_amount, current_amount) VALUES ('test category2', 100, 1, 2, 100, 100, 2, 200, 200);`

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
		t.Error("couldn't insert tags;", err)
		return err
	}

	return nil
}

func getUser(t *testing.T) (models.User, error) {
	// Get user
	query := `SELECT id, email, password, db_version, free_funds, created_at, updated_at FROM user;`

	// Get rows
	row := db.QueryRow(query)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.DBVersion,
		&user.FreeFunds,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		t.Error("couldn't get user", err)
		return models.User{}, err
	}

	return user, nil
}

func getAccounts(t *testing.T) ([]models.Account, error) {
	// Get accounts
	query := `SELECT id, name, current_amount, usage_count, table_order, created_at, updated_at FROM accounts;`

	// Get rows
	rows, err := db.Query(query)
	if err != nil {
		t.Error("couldn't get accounts", err)
		return nil, err
	}

	var accounts []models.Account

	for rows.Next() {
		var account models.Account

		err = rows.Scan(
			&account.ID,
			&account.Name,
			&account.CurrentAmount,
			&account.UsageCount,
			&account.TableOrder,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			t.Error("couldn't scan accounts", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func getCategories(t *testing.T) ([]models.Category, error) {
	// Get categories
	query := `SELECT
		id,
		name,
		budget_input,
		last_input_date,
		concat(input_interval, input_period) as input_interval,
		spending_limit,
		spending_left,
		initial_amount,
		current_amount,
		table_order,
		created_at,
		updated_at
	FROM categories;`

	// Get rows
	rows, err := db.Query(query)
	if err != nil {
		t.Error("couldn't get categories", err)
		return nil, err
	}

	var categories []models.Category

	for rows.Next() {
		var category models.Category

		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.BudgetInput,
			&category.LastInputDate,
			&category.InputInterval,
			&category.SpendingLimit,
			&category.SpendingLeft,
			&category.InitialAmount,
			&category.CurrentAmount,
			&category.TableOrder,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			t.Error("couldn't scan category", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func getExpenses(t *testing.T) ([]models.Expense, error) {
	// Get expenses
	query := `SELECT
		id,
		amount,
		date,
		from_account,
		from_category,
		created_at,
		updated_at
	FROM expenses;`

	// Get rows
	rows, err := db.Query(query)
	if err != nil {
		t.Error("couldn't get expenses", err)
		return nil, err
	}

	var expenses []models.Expense

	for rows.Next() {
		var expense models.Expense

		err = rows.Scan(
			&expense.ID,
			&expense.Amount,
			&expense.Date,
			&expense.FromAccountId,
			&expense.FromCategoryId,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			t.Error("couldn't scan category", err)
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func getTags(t *testing.T) ([]models.Tag, error) {
	// Get tags
	query := `SELECT
		id,
		name,
		usage_count,
		created_at,
		updated_at
	FROM tags;`

	// Get rows
	rows, err := db.Query(query)
	if err != nil {
		t.Error("couldn't get tags", err)
		return nil, err
	}

	var tags []models.Tag

	for rows.Next() {
		var tag models.Tag

		err = rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.UsageCount,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			t.Error("couldn't scan tag", err)
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func getExpenseTags(t *testing.T) ([]models.ExpenseToTagRealtion, error) {
	// Get expense tags
	query := `SELECT
		id,
		expense_id,
		tag_id,
		created_at,
		updated_at
	FROM expense_tags;`

	// Get rows
	rows, err := db.Query(query)
	if err != nil {
		t.Error("couldn't get expense tags", err)
		return nil, err
	}

	var expenseTags []models.ExpenseToTagRealtion

	for rows.Next() {
		var expenseTag models.ExpenseToTagRealtion

		err = rows.Scan(
			&expenseTag.ID,
			&expenseTag.ExpenseId,
			&expenseTag.TagId,
			&expenseTag.CreatedAt,
			&expenseTag.UpdatedAt,
		)
		if err != nil {
			t.Error("couldn't scan tag", err)
			return nil, err
		}
		expenseTags = append(expenseTags, expenseTag)
	}

	return expenseTags, nil
}
