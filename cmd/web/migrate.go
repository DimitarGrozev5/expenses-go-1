package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
)

func Migrate(dbName string) error {

	// Delete old DB
	os.Remove(dbName)

	_, err := os.OpenFile(dbName, os.O_RDONLY, os.ModeType)
	if err == nil {
		log.Printf("DB not deleted!!!!!!!!!!!")
		return err
	}

	// Create db
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s.db?_fk=%s", dbName, url.QueryEscape("true")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*
	 * User table
	 *
	 * Each user has his own database
	 * The user table contains the user name and password
	 * It contains the current DB version. This will become important when the app is live, so it will be able to migrate the user databases when an update is pushed
	 * In the future it can contain user settings if the need arises
	 */
	stmt := `CREATE TABLE user (
					id			INTEGER					NOT NULL	PRIMARY KEY		AUTOINCREMENT,

					email		TEXT		UNIQUE		NOT NULL,
					password	TEXT					NOT NULL,
					db_version	INTEGER

					created_at	DATETIME				NOT NULL	DEFAULT CURRENT_TIMESTAMP,
					updated_at	DATETIME							DEFAULT null
				)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Expenses table
	 *
	 * Contains data about user expenses
	 * Amount can be both positive or negative
	 * Date is the date the expense was made
	 * It tracks from which account was the money taken
	 * It tracks from which category was the money taken
	 */
	stmt = `CREATE TABLE expenses (
		id				INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

		amount			NUMERIC		NOT NULL,
		date			DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,

		from_account	INTEGER		NOT NULL	REFERENCES accounts (id)
													ON DELETE RESTRICT,
		from_category	INTEGER		NOT NULL	REFERENCES categories (id)
													ON DELETE RESTRICT,
		
		created_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at		DATETIME				DEFAULT null
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Tags table
	 *
	 * Expenses can be tagged
	 * Each expense is meant to have at least one tag
	 * The tags table tracks how many times is each tag used and when was the last useage time
	 */
	stmt = `CREATE TABLE tags (
		id			INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

		name		TEXT		NOT NULL	UNIQUE,
		usage_count	INTEGER		NOT NULL					DEFAULT 0,

		created_at	DATETIME	NOT NULL					DEFAULT CURRENT_TIMESTAMP,
		updated_at	DATETIME								DEFAULT null
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	// Many-to-many relationa table to link expenses and tags
	stmt = `CREATE TABLE expense_tags (
		id			INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

		expense_id	INTEGER		NOT NULL	REFERENCES expenses (id)
												ON DELETE CASCADE
												ON UPDATE CASCADE,
		tag_id		INTEGER		NOT NULL	REFERENCES tags (id)
												ON DELETE CASCADE
												ON UPDATE CASCADE,

		created_at	DATETIME	NOT NULL					DEFAULT CURRENT_TIMESTAMP,
		updated_at	DATETIME								DEFAULT null
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Triggers related to expenses and tags
	 *
	 * When a tag is linked to an expense
	 * or a relation is changed
	 * or a relation is removed
	 * Update the affected tag usage count
	 * Don't allow deletion of tags that have an usage count greather than zero
	 */
	stmt = `	CREATE TRIGGER tag_usage_count_insert
					AFTER INSERT
					ON expense_tags
				BEGIN
					UPDATE tags SET
						usage_count = usage_count + 1,
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
				END;
				
				CREATE TRIGGER tags_prevent_deletion
					BEFORE DELETE
					ON tags
					WHEN old.usage_count > 0
				BEGIN
					SELECT RAISE (ABORT, 'cant delete tags that are used');
				END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	// TODO: make default values of inital_amount 0; It's 100 for dev purposes
	/*
	 * Accounts table
	 *
	 * The account is meant to represent an account that can be used to store money
	 * This can mean a bank account, a piggy bank, or even the cash you keep on hand
	 *
	 * The account keeps track of it's money by having an initial amount
	 * Every expense that is taken from the account, removes money from the inital amount
	 * The current amount is stored in current_amount, to make the value easily available and is calculated whenever and amount is moved from the account
	 * Money is added to the account through the categories. Whenever money is added to a category, the initial amount and the current amount are updated
	 *
	 * The account keeps a record of how often it is used - usage count, so the frontend can place the more used accounts infront in selects
	 * TODO: maybe implement some sort of filter, to smooth the usage count, so account order doesn't change all the time
	 *
	 * The account keep a table order, so the frontend can show the accounts in the order the user prefers. The order is in DESC order
	 * The initial order is set automatically when the account is inserted
	 */
	stmt = `CREATE TABLE accounts (
		id				INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,
		
		name			TEXT		NOT NULL	UNIQUE		CHECK (length(name) > 3),
		initial_amount	NUMERIC		NOT NULL	DEFAULT 0	CHECK (initial_amount >= 0),
		current_amount	NUMERIC		NOT NULL	DEFAULT 0	CHECK (current_amount >= 0),

		usage_count		INTEGER		NOT NULL	DEFAULT 0,
		table_order		INTEGER		NOT NULL	DEFAULT -1,

		created_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at		DATETIME				DEFAULT null
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Account related triggers
	 *
	 * When an expense is added
	 * Update the related account current_amount and usage_count
	 *
	 * When an account is created, update it's table_order so it is the next available number
	 * Prevent the order being changed bellow 0 and above the current max number
	 * When an account is deleted, update all of the table orders so there are no gaps
	 *
	 * Don't allow the deletion of accounts that have expenses linked to them
	 */
	stmt = `    CREATE TRIGGER accounts_update_current_amount_when_initial_amount_changes
					AFTER UPDATE
					ON accounts
					WHEN old.initial_amount <> new.initial_amount
				BEGIN
					UPDATE accounts SET 
	 					current_amount = current_amount - old.initial_amount + new.initial_amount,
						updated_at = datetime('now')
					WHERE accounts.id = old.id;
				END;
				
				CREATE TRIGGER account_current_amount_add
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

				CREATE TRIGGER accounts_check_table_order
					BEFORE UPDATE
					ON accounts
					WHEN old.table_order <> new.table_order
				BEGIN
					SELECT
						CASE
							WHEN new.table_order < 1 THEN
								RAISE (ABORT, 'Cant move the first account up')
							WHEN new.table_order > (SELECT COUNT(*) from accounts) THEN
								RAISE (ABORT, 'Cant move the last account down')
						END;
				END;

				CREATE TRIGGER accounts_auto_update_account_order
					BEFORE UPDATE
					ON accounts
					WHEN old.table_order <> new.table_order
				BEGIN
					UPDATE accounts SET
						table_order = old.table_order
					WHERE accounts.table_order = new.table_order;
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
		return err
	}

	/*
	 * Categories table
	 *
	 *
	 * The Budget of the Category is constucted by specifying a couple of parameters:
	 *
	 * Budget input - the amount the user expects to input in the budget
	 * Last input date - the last time money was placed in the category
	 * Input interval period - the period over which the user expect to input to the budget - can be YEARS, MONTHS, DAYS
	 * Input interval value - the amount of periods. e.g. 3 DAYS, 4 MONTHS
	 *
	 * Spending limit - the maximum amount the user plans to spend
	 * Last spending reset - the last time the spending limit was reset
	 * Spending interval period - the period over which the spending limit has been imposed
	 * Spending interval value - the amount of periods. e.g. 3 DAYS, 1 MONTHS
	 *
	 *
	 * The category has an initial amount that defaults to zero
	 * When the user adds funds to the category, the initial amount is updated
	 * The category has a current amount that tracks the current available funds
	 * When the user adds expenses the current amount is updated
	 *
	 * The category keep a table order, so the frontend can show the categories in the order the user prefers. The order is in DESC
	 * The initial order is set automatically when the caegory is inserted
	 */
	stmt = ` CREATE TABLE categories (
		id						INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,
		
		name					TEXT		NOT NULL	UNIQUE,

		budget_input			NUMERIC		NOT NULL		CHECK (budget_input >= 0),
		last_input_date			DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		input_interval			INTEGER		NOT NULL		CHECK (input_interval > 0),
		input_period			INTEGER		NOT NULL	REFERENCES time_periods (id)
															ON DELETE RESTRICT,

		spending_limit			NUMERIC		NOT NULL		CHECK (budget_input >= 0),
		last_spending_reset		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		spending_interval		INTEGER		NOT NULL		CHECK (input_interval > 0),
		spending_period			INTEGER		NOT NULL	REFERENCES time_periods (id)
															ON DELETE RESTRICT,

		initial_amount			NUMERIC		NOT NULL	DEFAULT 0		CHECK (initial_amount >= 0),
		current_amount			NUMERIC		NOT NULL	DEFAULT 0		CHECK (current_amount >= 0),

		table_order				INTEGER		NOT NULL	DEFAULT -1,

		created_at				DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at				DATETIME	null		DEFAULT null
	)
	`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Timeframes enum table
	 *
	 * Create table with very restricted values
	 *
	 * Add trigger to prevent deleting values
	 *
	 * Seed values
	 */
	stmt = `	CREATE TABLE time_periods (
		id						INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

		period					TEXT		NOT NULL	UNIQUE			CHECK (period IN (' YEARS', ' MONTHS', ' DAYS')),

		created_at				DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at				DATETIME				DEFAULT null
	 );

	 CREATE TRIGGER dont_delete_from_time_periods
		 BEFORE DELETE
		 ON time_periods
	 BEGIN
		 SELECT RAISE (ABORT, 'Cant delete from time_periods');
	 END;

	 CREATE TRIGGER dont_update_from_time_periods
		 BEFORE UPDATE
		 ON time_periods
	 BEGIN
		 SELECT RAISE (ABORT, 'Cant update in time_periods');
	 END;

	 INSERT INTO time_periods (period) VALUES (' YEARS');
	 INSERT INTO time_periods (period) VALUES (' MONTHS');
	 INSERT INTO time_periods (period) VALUES (' DAYS');
	 `

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Triggers related to categories
	 *
	 * When user updates category initial_amount, update current amount value
	 *
	 * When expense is changed, update current amount
	 *
	 * Handle table order changes similar to accounts
	 *
	 * Block deletion of accounts with funds in them
	 */
	stmt = `
	 			CREATE TRIGGER categories_update_current_amount_when_initial_amount_changes
					AFTER UPDATE
					ON categories
					WHEN old.initial_amount <> new.initial_amount
				BEGIN
					UPDATE categories SET 
	 					current_amount = current_amount - old.initial_amount + new.initial_amount,
						updated_at = datetime('now')
					WHERE categories.id = old.id;
				END;



				CREATE TRIGGER categories_expenses_is_added
					BEFORE INSERT
					ON expenses
				BEGIN
					UPDATE categories SET
						current_amount = current_amount - new.amount,
						updated_at = datetime('now')
					WHERE categories.id = new.from_category;
				END;
				
				CREATE TRIGGER categories_expense_is_removed
					BEFORE DELETE
					ON expenses
				BEGIN
					UPDATE categories SET
						current_amount = current_amount + old.amount,
						updated_at = datetime('now')
					WHERE categories.id = old.from_category;
				END;
				
				CREATE TRIGGER categories_expense_amount_is_changed
					BEFORE UPDATE
					ON expenses
					WHEN
						old.amount <> new.amount AND
						old.from_category = new.from_category
				BEGIN
					UPDATE categories SET
						current_amount = current_amount + old.amount - new.amount,
						updated_at = datetime('now')
					WHERE categories.id = new.from_category;
				END;

				CREATE TRIGGER categories_expense_category_is_changed
					BEFORE UPDATE
					ON expenses
					WHEN old.from_category <> new.from_category
				BEGIN
					UPDATE categories SET
						current_amount = current_amount + old.amount,
						usage_count = usage_count - 1
					WHERE categories.id = old.from_category;
					
					UPDATE categories SET
						current_amount = current_amount - new.amount,
						usage_count = usage_count + 1,
						updated_at = datetime('now')
					WHERE categories.id = new.from_category;
				END;




				CREATE TRIGGER categories_set_order_for_new_accounts
					AFTER INSERT
					ON categories
				BEGIN
					UPDATE categories SET
						table_order = (SELECT COUNT(*) FROM categories),
						updated_at = datetime('now')
					WHERE categories.table_order = -1;
				END;

				CREATE TRIGGER categories_stop_invalid_table_order_values
					BEFORE UPDATE
					ON categories
					WHEN old.table_order <> new.table_order
				BEGIN
					SELECT
						CASE
							WHEN new.table_order < 1 THEN
								RAISE (ABORT, 'Cant move the first category up')
							WHEN new.table_order > (SELECT COUNT(*) from categories) THEN
								RAISE (ABORT, 'Cant move the last category down')
						END;
				END;

				CREATE TRIGGER categories_auto_update_account_order
					BEFORE UPDATE
					ON categories
					WHEN old.table_order <> new.table_order
				BEGIN
					UPDATE categories SET
						table_order = old.table_order
					WHERE categories.table_order = new.table_order;
				END;

				CREATE TRIGGER categories_update_order_after_delete
					AFTER DELETE
					ON categories
				BEGIN
					UPDATE categories SET
						table_order = table_order - 1
					WHERE categories.table_order > old.table_order;
				END;


				

				CREATE TRIGGER categories_block_delete
					BEFORE DELETE
					ON categories
				BEGIN
					SELECT
						CASE
							WHEN old.current_amount > 0 THEN
								RAISE (ABORT, 'Cant delete category that is being used')
						END;
				END;
	`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	return nil
}
