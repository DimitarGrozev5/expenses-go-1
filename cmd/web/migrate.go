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
	stmt := `	CREATE TABLE user (
					id			INTEGER					NOT NULL	PRIMARY KEY		AUTOINCREMENT,

					email		TEXT		UNIQUE		NOT NULL,
					password	TEXT					NOT NULL,
					db_version	INTEGER,
					free_funds	NUMERIC					NOT NULL	DEFAULT 0	CHECK (free_funds >= 0),

					created_at	DATETIME				NOT NULL	DEFAULT CURRENT_TIMESTAMP,
					updated_at	DATETIME							DEFAULT null
				);
				
				CREATE VIEW procedure_add_free_funds AS
					SELECT free_funds as amount, null as to_account FROM user;
					
				CREATE TRIGGER triggers__procedure_add_free_funds__update
					INSTEAD OF INSERT
					ON procedure_add_free_funds
				BEGIN
					UPDATE user SET
						free_funds = free_funds + new.amount;

					UPDATE accounts SET
						current_amount = current_amount + new.amount
					WHERE id = new.to_account;
				END;`

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

	/**
	 **
	 ** Expenses procedures
	 **
	 **/

	/*
	 * Insert new expense
	 */
	stmt = `CREATE VIEW procedure_new_expense AS
				SELECT amount, date, from_account, from_category FROM expenses;
				
			CREATE TRIGGER triggers__procedure_new_expense__add_expense
				INSTEAD OF INSERT
				ON procedure_new_expense
			BEGIN
				INSERT INTO expenses (
					amount,
					date,
					from_account,
					from_category
				) VALUES (
					new.amount,
					new.date,
					new.from_account,
					new.from_category
				);

				UPDATE accounts SET
					current_amount = current_amount - new.amount,
					usage_count = usage_count + 1,
					updated_at = datetime('now')
				WHERE accounts.id = new.from_account;

				UPDATE categories SET
					current_amount = current_amount - new.amount,
					updated_at = datetime('now')
				WHERE categories.id = new.from_category;
			END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Delete expense
	 */
	stmt = `CREATE VIEW procedure_remove_expense AS
				SELECT id, amount, from_account, from_category FROM expenses;
			
			CREATE TRIGGER triggers__procedure_remove_expense_remove
				INSTEAD OF DELETE
				ON procedure_remove_expense
			BEGIN
				DELETE FROM expenses WHERE id = old.id;

				UPDATE accounts SET
					current_amount = current_amount + old.amount,
					usage_count = usage_count - 1,
					updated_at = datetime('now')
				WHERE accounts.id = old.from_account;
				
				UPDATE categories SET
					current_amount = current_amount + old.amount,
					updated_at = datetime('now')
				WHERE categories.id = old.from_category;
			END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Update expense
	 */
	stmt = `CREATE VIEW procedure_update_expense AS
				SELECT id, amount, date, from_account, from_category FROM expenses;
			
			CREATE TRIGGER triggers__procedure_update_expense__update
				INSTEAD OF UPDATE
				ON procedure_update_expense
			BEGIN
				UPDATE expenses SET
					amount = 		COALESCE(new.amount, old.amount),
					date = 			COALESCE(new.date, old.date),
					from_account = 	COALESCE(new.from_account, old.from_account),
					from_category =	COALESCE(new.from_category, old.from_category),
					updated_at = datetime('now')
				WHERE id = new.id;
			END;
			
			CREATE TRIGGER triggers__procedure_update_expense__amount_changes_account_stays_the_same
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
			
			CREATE TRIGGER triggers__procedure_update_expense__account_changes
				BEFORE UPDATE
				ON expenses
				WHEN old.from_account <> new.from_account
			BEGIN
				UPDATE accounts SET
					current_amount = current_amount + old.amount,
					usage_count = usage_count - 1,
					updated_at = datetime('now')
				WHERE accounts.id = old.from_account;
				
				UPDATE accounts SET
					current_amount = current_amount - new.amount,
					usage_count = usage_count + 1,
					updated_at = datetime('now')
				WHERE accounts.id = new.from_account;
			END;
			
			CREATE TRIGGER triggers__procedure_update_expense__amount_changes_category_stays_the_same
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

			CREATE TRIGGER triggers__procedure_update_expense__category_changes
				BEFORE UPDATE
				ON expenses
				WHEN old.from_category <> new.from_category
			BEGIN
				UPDATE categories SET
					current_amount = current_amount + old.amount,
					updated_at = datetime('now')
				WHERE categories.id = old.from_category;
				
				UPDATE categories SET
					current_amount = current_amount - new.amount,
					updated_at = datetime('now')
				WHERE categories.id = new.from_category;
			END;`

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
		);
	
		CREATE VIEW procedure_insert_tag AS
			SELECT name FROM tags;
		
		CREATE TRIGGER trigger__procedure_insert_tag__insert
			INSTEAD OF INSERT
			ON procedure_insert_tag
		BEGIN
			INSERT INTO tags (name) VALUES (new.name)
				ON CONFLICT (name) DO NOTHING;
		END;
		
		CREATE VIEW procedure_remove_tag AS
				SELECT id, usage_count FROM tags;
			
		CREATE TRIGGER triggers__procedure_remove_tag
			INSTEAD OF DELETE
			ON procedure_remove_tag
		BEGIN
			DELETE FROM tags WHERE id = old.id;
		END;`

	// CREATE TRIGGER triggers__expense_tags__tags_prevent_deletion
	// 				BEFORE DELETE
	// 				ON tags
	// 				WHEN old.usage_count > 0
	// 			BEGIN
	// 				SELECT RAISE (ABORT, 'cant delete tags that are used');
	// 			END;

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	// Many-to-many relations table to link expenses and tags
	stmt = `CREATE TABLE expense_tags (
				id			INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

				expense_id	INTEGER		NOT NULL	REFERENCES expenses (id)
														ON DELETE CASCADE
														ON UPDATE CASCADE,
				tag_id		INTEGER		NOT NULL	REFERENCES tags (id)
														ON DELETE RESTRICT
														ON UPDATE CASCADE,

				created_at	DATETIME	NOT NULL					DEFAULT CURRENT_TIMESTAMP,
				updated_at	DATETIME								DEFAULT null
			);

			
			CREATE VIEW procedure_link_tag_to_expense AS
				SELECT expense_id, tag_id FROM expense_tags;
				
			CREATE TRIGGER trigger__procedure_link_tag_to_expense__add
				INSTEAD OF INSERT
				ON procedure_link_tag_to_expense
			BEGIN
				INSERT INTO expense_tags (expense_id, tag_id) VALUES (new.expense_id, new.tag_id);
			END;
			
			
			CREATE VIEW procedure_unlink_tag_from_expense AS
				SELECT id, expense_id, tag_id FROM expense_tags;
				
			CREATE TRIGGER trigger__procedure_unlink_tag_from_expense__remove
				INSTEAD OF DELETE
				ON procedure_unlink_tag_from_expense
			BEGIN
				DELETE FROM expense_tags WHERE id = old.id;
			END;`

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
	stmt = `	CREATE TRIGGER triggers__expense_tags__tag_usage_count_insert
					AFTER INSERT
					ON expense_tags
				BEGIN
					UPDATE tags SET
						usage_count = usage_count + 1,
						updated_at = datetime('now')
					WHERE tags.id = new.tag_id;
				END;
				
				CREATE TRIGGER triggers__expense_tags__tag_usage_count_update
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
				
				CREATE TRIGGER triggers__expense_tags__tag_usage_count_delete
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
		return err
	}

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
		current_amount	NUMERIC		NOT NULL	DEFAULT 0	CHECK (current_amount >= 0),

		usage_count		INTEGER		NOT NULL	DEFAULT 0,
		table_order		INTEGER		NOT NULL				CHECK (table_order > 0),

		created_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at		DATETIME				DEFAULT null
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/**
	 **
	 ** Account procedures
	 **
	 **/

	/*
	 * Create new account
	 *
	 * Provide only a name. Everything else is auto generated
	 */
	stmt = `CREATE VIEW procedure_new_account AS
				SELECT name FROM accounts;
				
			CREATE TRIGGER trigger__procedure_new_account__insert_new_account
				INSTEAD OF INSERT
				ON procedure_new_account
			BEGIN
				INSERT INTO accounts (
					name,
					table_order
				) VALUES (
					new.name,
					(SELECT COUNT(*) FROM accounts) + 1
				);
			END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Edit account name
	 */
	stmt = `CREATE VIEW procedure_account_name AS
				SELECT id, name FROM accounts;
				
			CREATE TRIGGER trigger__procedure_account_name__update_account_name
				INSTEAD OF UPDATE
				ON procedure_account_name
				WHEN
					old.name <> new.name AND
					old.id = new.id
			BEGIN
				UPDATE accounts SET name=new.name WHERE id=new.id;
			END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Change accounts order
	 *
	 * When updating, pass new table_order value
	 * Too big or too small values throw an error
	 */
	stmt = `CREATE VIEW procedure_change_accounts_order AS
	 			SELECT id, table_order FROM accounts;
				
			CREATE TRIGGER trigger__procedure_change_accounts_order__swap_table_orders
				INSTEAD OF UPDATE
				ON procedure_change_accounts_order
				WHEN
					old.table_order <> new.table_order AND
					old.id = new.id
			BEGIN
				SELECT
					CASE
						WHEN new.table_order < 1 THEN
							RAISE (ABORT, 'cant move last account down')
						WHEN new.table_order > (SELECT COUNT(*) from accounts) THEN
							RAISE (ABORT, 'cant move first account up')
					END;
				
				UPDATE accounts SET
					table_order = old.table_order
				WHERE table_order = new.table_order;

				UPDATE accounts SET
					table_order = new.table_order
				WHERE id = new.id;
			END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Delete account
	 *
	 * Remove account and update table order of the rest of the accounts so there are no gaps
	 */
	stmt = `CREATE VIEW procedure_remove_account AS
				SELECT id, table_order FROM accounts;
			
			CREATE TRIGGER triggers__procedure_remove_account__delete_account
				INSTEAD OF DELETE
				ON procedure_remove_account
			BEGIN
				SELECT
					CASE
						WHEN current_amount > 0 THEN
							RAISE (ABORT, 'cant delete an account that is used')
					END
				FROM accounts
				WHERE id = old.id;

				DELETE FROM accounts WHERE id=old.id;

				UPDATE accounts SET
					table_order = table_order - 1
				WHERE table_order > old.table_order;
			END;`

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
	 * Spending left - the amount of money the user is allowed to spend until the end of the spending period
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
		spending_left			NUMERIC		NOT NULL,
		last_spending_reset		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		spending_interval		INTEGER		NOT NULL		CHECK (input_interval > 0),
		spending_period			INTEGER		NOT NULL	REFERENCES time_periods (id)
															ON DELETE RESTRICT,

		initial_amount			NUMERIC		NOT NULL	DEFAULT 0		CHECK (initial_amount >= 0),
		current_amount			NUMERIC		NOT NULL	DEFAULT 0		CHECK (current_amount >= 0),

		table_order				INTEGER		NOT NULL					CHECK (table_order > 0),

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

	/**
	 **
	 ** Category stored procedures
	 **
	 **/

	/*
	 * Create new category
	 */
	stmt = `CREATE VIEW procedure_new_category AS
				SELECT 
					name,
					budget_input,
					input_interval,
					input_period,
					spending_limit,
					spending_interval,
					spending_period
				FROM categories;
				
			CREATE TRIGGER triggers__procedure_new_category__insert_new
				INSTEAD OF INSERT
				ON procedure_new_category
			BEGIN
				INSERT INTO categories (
					name,
					budget_input,
					input_interval,
					input_period,
					spending_limit,
					spending_left,
					spending_interval,
					spending_period,
					initial_amount,
					current_amount,
					table_order
				) VALUES (
					new.name,
					new.budget_input,
					new.input_interval,
					new.input_period,
					new.spending_limit,
					new.spending_limit,
					new.spending_interval,
					new.spending_period,
					0,
					0,
					(SELECT COUNT(*) FROM categories) + 1
				);
			END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Update category name
	 */
	stmt = `CREATE VIEW procedure_category_name AS
				SELECT id, name FROM categories;
				
			CREATE TRIGGER trigger__procedure_category_name__update_name
				INSTEAD OF UPDATE
				ON procedure_category_name
				WHEN
					old.name <> new.name AND
					old.id = new.id
				BEGIN
					UPDATE categories SET
						name = new.name
					WHERE id = new.id;
				END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Update categories order
	 *
	 * When updating, pass new table_order value
	 * Too big or too small values throw an error
	 */
	stmt = `CREATE VIEW procedure_change_categories_order AS
	 			SELECT id, table_order FROM categories;
				
			CREATE TRIGGER trigger__procedure_change_categories_order__swap_table_order
				INSTEAD OF UPDATE
				ON procedure_change_categories_order
				WHEN
					old.table_order <> new.table_order AND
					old.id = new.id
			BEGIN
				SELECT
					CASE
						WHEN new.table_order < 1 THEN
							RAISE (ABORT, 'cant move last category down')
						WHEN new.table_order > (SELECT COUNT(*) from categories) THEN
							RAISE (ABORT, 'cant move first category up')
					END;
				
				UPDATE categories SET
					table_order = old.table_order
				WHERE table_order = new.table_order;

				UPDATE categories SET
					table_order = new.table_order
				WHERE id = new.id;
			END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Procedure to delete unused categories
	 *
	 * A category is unused when it's not referenced by an other table and doens't have funds
	 */
	stmt = `CREATE VIEW procedure_remove_category AS
				SELECT id, table_order FROM categories;
				
			CREATE TRIGGER trigger__procedure_remove_category__remove_category
				INSTEAD OF DELETE
				ON procedure_remove_category
			BEGIN
				SELECT
					CASE
						WHEN initial_amount > 0 THEN
							RAISE (ABORT, 'cant delete a category that is used')
					END
				FROM categories
				WHERE id = old.id;

				DELETE FROM categories WHERE id = old.id;

				UPDATE categories SET
					table_order = table_order - 1
				WHERE table_order > old.table_order;
			END;`

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

	return nil
}
