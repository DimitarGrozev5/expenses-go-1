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
						free_funds = free_funds + new.amount,
						updated_at = datetime('now');

					UPDATE accounts SET
						current_amount = current_amount + new.amount,
						updated_at = datetime('now')
					WHERE id = new.to_account;

					INSERT INTO accounts_input_log (
						account,
						amount
					) VALUES (
						new.to_account,
						new.amount
					);
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
	 *
	 * From period tracks if the expense is archived
	 * If the expense is not archived, the field is null
	 * When the expense gets archived it will reference the periods table
	 */
	stmt = `CREATE TABLE expenses (
		id				INTEGER		NOT NULL		PRIMARY KEY		AUTOINCREMENT,

		amount			NUMERIC		NOT NULL,
		date			DATETIME	NOT NULL		DEFAULT CURRENT_TIMESTAMP,

		from_account	INTEGER		NOT NULL		REFERENCES accounts (id)
														ON UPDATE CASCADE
														ON DELETE RESTRICT,
		from_category	INTEGER		NOT NULL		REFERENCES categories (id)
														ON UPDATE CASCADE
														ON DELETE RESTRICT,

		from_period		INTEGER		DEFAULT null	REFERENCES archived_periods (id)
														ON UPDATE CASCADE
														ON DELETE RESTRICT,
		
		created_at		DATETIME	NOT NULL		DEFAULT CURRENT_TIMESTAMP,
		updated_at		DATETIME					DEFAULT null
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * View current expenses
	 */
	stmt = `CREATE VIEW view_current_expenses AS
				SELECT id, amount, date, from_account, from_category, created_at, updated_at
				FROM expenses
				WHERE from_period = null;`

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
					spending_left = spending_left - new.amount,
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
				SELECT id, amount, from_account, from_category, from_period FROM expenses;
			
			CREATE TRIGGER triggers__procedure_remove_expense__remove
				INSTEAD OF DELETE
				ON procedure_remove_expense
			BEGIN
				SELECT
					CASE
						WHEN old.from_period IS NOT NULL THEN
							RAISE (ABORT, 'cant delete archived expense')
					END;

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
			END;
			`

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
				SELECT id, amount, date, from_account, from_category, from_period FROM expenses;
			
			CREATE TRIGGER triggers__procedure_update_expense__update
				INSTEAD OF UPDATE
				ON procedure_update_expense
			BEGIN
			SELECT
				CASE
					WHEN old.from_period IS NOT NULL THEN
						RAISE (ABORT, 'cant delete archived expense')
				END;
				
				UPDATE expenses SET
					amount = 		COALESCE(new.amount, old.amount),
					date = 			COALESCE(new.date, old.date),
					from_account = 	COALESCE(new.from_account, old.from_account),
					from_category =	COALESCE(new.from_category, old.from_category),
					updated_at = datetime('now')
				WHERE id = old.id;
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

	/*
	 * Accounts input log
	 *
	 * When account current_amount increases, log changes
	 */
	stmt = `CREATE TABLE accounts_input_log (
				id				INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

				account			INTEGER		NOT NULL	REFERENCES accounts (id)
															ON UPDATE CASCADE
															ON DELETE CASCADE,
				amount			NUMERIC		NOT NULL,

				created_at		DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
				updated_at		DATETIME				DEFAULT null
			);`

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
	stmt = `CREATE VIEW procedure_insert_account AS
				SELECT name FROM accounts;
				
			CREATE TRIGGER trigger__procedure_insert_account__insert
				INSTEAD OF INSERT
				ON procedure_insert_account
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
	stmt = `CREATE VIEW procedure_account_update_name AS
				SELECT id, name FROM accounts;
				
			CREATE TRIGGER trigger__procedure_account_update_name__update
				INSTEAD OF UPDATE
				ON procedure_account_update_name
				WHEN old.name <> new.name
			BEGIN
				UPDATE accounts SET
					name=new.name,
					updated_at = datetime('now')
				WHERE id=old.id;
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
				WHEN old.table_order <> new.table_order
			BEGIN
				SELECT
					CASE
						WHEN new.table_order < 1 THEN
							RAISE (ABORT, 'cant move last account down')
						WHEN new.table_order > (SELECT COUNT(*) from accounts) THEN
							RAISE (ABORT, 'cant move first account up')
					END;
				
				UPDATE accounts SET
					table_order = old.table_order,
					updated_at = datetime('now')
				WHERE table_order = new.table_order;

				UPDATE accounts SET
					table_order = new.table_order,
					updated_at = datetime('now')
				WHERE id = old.id;
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
					table_order = table_order - 1,
					updated_at = datetime('now')
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
															ON UPDATE CASCADE
															ON DELETE RESTRICT,

		spending_limit			NUMERIC		NOT NULL		CHECK (spending_limit >= 0),
		spending_left			NUMERIC		NOT NULL		CHECK (spending_left >= 0),

		initial_amount			NUMERIC		NOT NULL	DEFAULT 0		CHECK (initial_amount >= 0),
		current_amount			NUMERIC		NOT NULL	DEFAULT 0		CHECK (current_amount >= 0),

		table_order				INTEGER		NOT NULL					CHECK (table_order > 0),

		created_at				DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at				DATETIME	null		DEFAULT null
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * View categories
	 */
	stmt = `CREATE VIEW view_categories AS
				SELECT
					c.id,
					c.name,
					c.budget_input,
					c.last_input_date,
					concat(c.input_interval, p.period) as input_interval,
					c.spending_limit,
					c.spending_left,
					c.initial_amount,
					c.current_amount,
					c.table_order,
					c.created_at,
					c.updated_at
				FROM categories AS c
				JOIN time_periods AS p ON c.input_period = p.id;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * View categories overview
	 */
	stmt = `CREATE VIEW view_categories_overview AS
				SELECT
					c.id,
					c.name,
					c.budget_input,
					c.input_interval,
					c.input_period,
					c.spending_limit,
					c.spending_left,
					c.last_input_date AS period_start,
					datetime(c.last_input_date, concat(c.input_interval, p.period)) AS period_end,
					c.initial_amount,
					c.current_amount,
					((SELECT COUNT(*) FROM archived_periods WHERE category=c.id) = 0) AS can_be_deleted,
					c.table_order
				FROM categories AS c
				JOIN time_periods AS p ON c.input_period = p.id;`

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
					spending_limit
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
				WHEN old.name <> new.name
				BEGIN
					UPDATE categories SET
						name = new.name,
						updated_at = datetime('now')
					WHERE id = old.id;
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
				WHEN old.table_order <> new.table_order
			BEGIN
				SELECT
					CASE
						WHEN new.table_order < 1 THEN
							RAISE (ABORT, 'cant move last category down')
						WHEN new.table_order > (SELECT COUNT(*) from categories) THEN
							RAISE (ABORT, 'cant move first category up')
					END;
				
				UPDATE categories SET
					table_order = old.table_order,
					updated_at = datetime('now')
				WHERE table_order = new.table_order;

				UPDATE categories SET
					table_order = new.table_order,
					updated_at = datetime('now')
				WHERE id = old.id;
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
						WHEN (SELECT COUNT(*) FROM archived_periods WHERE category = old.id) > 0 THEN
							RAISE (ABORT, 'cant delete a category that is used')
					END
				FROM categories
				WHERE id = old.id;

				DELETE FROM categories WHERE id = old.id;

				UPDATE categories SET
					table_order = table_order - 1,
					updated_at = datetime('now')
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
		caption					TEXT		NOT NULL,

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

	 INSERT INTO time_periods (period, caption) VALUES (' YEARS', 'Years');
	 INSERT INTO time_periods (period, caption) VALUES (' MONTHS', 'Months');
	 INSERT INTO time_periods (period, caption) VALUES (' DAYS', 'Days');
	 `

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Archived periods table
	 */
	stmt = `CREATE TABLE archived_periods (
		id						INTEGER		NOT NULL	PRIMARY KEY		AUTOINCREMENT,

		category				INTEGER		NOT NULL	REFERENCES categories (id)
															ON UPDATE CASCADE
															ON DELETE RESTRICT,
		period_start			DATETIME	NOT NULL,
		period_end				DATETIME	NOT NULL,

		budget_input			NUMERIC		NOT NULL,
		spending_limit			NUMERIC		NOT NULL,
		input_interval			INTEGER		NOT NULL,
		input_period			INTEGER		NOT NULL	REFERENCES time_periods (id)
															ON UPDATE CASCADE
															ON DELETE RESTRICT,

		initial_amount			NUMERIC		NOT NULL,
		end_amount				NUMERIC		NOT NULL,

		created_at				DATETIME	NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		updated_at				DATETIME				DEFAULT null
	)`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 *
	 * Input data to category and reset period
	 *
	 *
	 */
	stmt = `CREATE VIEW procedure_fund_category_and_reset_period AS
				SELECT
					current_amount as amount,
					id as category,
					budget_input,
					input_interval,
					input_period,
					spending_limit
				FROM categories;
				
			CREATE TRIGGER trigger__procedure_fund_category_and_reset_period_insert
				INSTEAD OF INSERT
				ON procedure_fund_category_and_reset_period
			BEGIN
				SELECT
					CASE
						WHEN new.amount < 0 THEN
							RAISE (ABORT, 'input amount must be greather than zero')
					END;

				INSERT INTO archived_periods (
					category,
					period_start,
					period_end,
					budget_input,
					spending_limit,
					input_interval,
					input_period,
					initial_amount,
					end_amount
				) SELECT
					id,
					last_input_date,
					datetime('now'),
					budget_input,
					spending_limit,
					input_interval,
					input_period,
					initial_amount,
					current_amount
				FROM categories
				WHERE id = new.category;
				
				UPDATE user SET
					free_funds = free_funds - new.amount,
					updated_at = datetime('now');

				UPDATE categories SET
					initial_amount = current_amount + new.amount,
					current_amount = current_amount + new.amount,
					budget_input = new.budget_input,
					input_interval = new.input_interval,
					input_period = new.input_period,
					spending_limit = new.spending_limit,
					spending_left = new.spending_limit,
					updated_at = datetime('now')
				WHERE id = new.category;

				UPDATE expenses SET
					from_period = last_insert_rowid(),
					updated_at = datetime('now')
				WHERE from_period IS NULL;
			END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	/*
	 * Add money to category without reseting the period
	 */
	stmt = `CREATE VIEW procedure_fund_category AS
				SELECT current_amount as amount, id as category FROM categories;
				
			CREATE TRIGGER triggers__procedure_fund_category__insert
				INSTEAD OF INSERT
				ON procedure_fund_category
			BEGIN
				UPDATE user SET
					free_funds = free_funds - new.amount,
					updated_at = datetime('now');

				UPDATE categories SET
					initial_amount = current_amount + new.amount,
					current_amount = current_amount + new.amount,
					updated_at = datetime('now')
				WHERE id = new.category;
			END;`

	// Execute query
	_, err = db.Exec(stmt)
	if err != nil {
		log.Printf("%q: %s\n", err, stmt)
		return err
	}

	return nil
}
