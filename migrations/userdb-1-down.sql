/*
 * Disable foreign key constraints just in case
 */
PRAGMA foreign_keys = OFF;

/*
 * Users table
 */
DROP TABLE IF EXISTS users;

DROP VIEW IF EXISTS procedure_add_free_funds;

DROP TRIGGER IF EXISTS triggers__procedure_add_free_funds__update;

/*
 * Expenses table
 */
DROP TABLE IF EXISTS expenses;

DROP VIEW IF EXISTS view_current_expenses;

DROP VIEW IF EXISTS view_detailed_expenses;

DROP VIEW IF EXISTS procedure_new_expense;

DROP TRIGGER IF EXISTS triggers__procedure_new_expense__add_expense;

DROP VIEW IF EXISTS procedure_remove_expense;

DROP TRIGGER IF EXISTS triggers__procedure_remove_expense__remove;

DROP VIEW IF EXISTS procedure_update_expense;

DROP TRIGGER IF EXISTS triggers__procedure_update_expense__update;

DROP TRIGGER IF EXISTS triggers__procedure_update_expense__amount_changes_account_stays_the_same;

DROP TRIGGER IF EXISTS triggers__procedure_update_expense__account_changes;

DROP TRIGGER IF EXISTS triggers__procedure_update_expense__amount_changes_category_stays_the_same;

DROP TRIGGER IF EXISTS triggers__procedure_update_expense__category_changes;

/*
 * Tags
 */
DROP TABLE IF EXISTS tags;

DROP IF EXISTS VIEW procedure_insert_tag;

DROP IF EXISTS TRIGGER trigger__procedure_insert_tag__insert;

DROP IF EXISTS VIEW procedure_remove_tag;

DROP IF EXISTS TRIGGER triggers__procedure_remove_tag;

/*
 * Expense to tag
 */
DROP TABLE IF EXISTS expense_tags;

DROP VIEW IF EXISTS procedure_link_tag_to_expense;

DROP TRIGGER IF EXISTS trigger__procedure_link_tag_to_expense__add;

DROP VIEW IF EXISTS procedure_unlink_tags_from_expense;

DROP TRIGGER IF EXISTS trigger__procedure_unlink_tags_from_expense__remove;

/*
 * Triggers for expenses and tags
 */
DROP TRIGGER IF EXISTS triggers__expense_tags__tag_usage_count_insert;

DROP TRIGGER IF EXISTS triggers__expense_tags__tag_usage_count_update;

DROP TRIGGER IF EXISTS triggers__expense_tags__tag_usage_count_delete;

/*
 * Accounts
 */
DROP TABLE IF EXISTS accounts;

DROP TABLE IF EXISTS accounts_input_log;

DROP VIEW IF EXISTS procedure_insert_account;

DROP TRIGGER IF EXISTS trigger__procedure_insert_account__insert;

DROP VIEW IF EXISTS procedure_account_update_name;

DROP TRIGGER IF EXISTS trigger__procedure_account_update_name__update;

DROP VIEW IF EXISTS procedure_change_accounts_order;

DROP TRIGGER IF EXISTS trigger__procedure_change_accounts_order__swap_table_orders;

DROP VIEW IF EXISTS procedure_remove_account;

DROP TRIGGER IF EXISTS triggers__procedure_remove_account__delete_account;

/*
 * Categories
 */
DROP TABLE IF EXISTS categories;

DROP VIEW IF EXISTS view_categories;

DROP VIEW IF EXISTS view_categories_overview;

DROP VIEW IF EXISTS procedure_new_category;

DROP TRIGGER IF EXISTS triggers__procedure_new_category__insert_new;

DROP VIEW IF EXISTS procedure_category_name;

DROP TRIGGER IF EXISTS trigger__procedure_category_name__update_name;

DROP VIEW IF EXISTS procedure_change_categories_order;

DROP TRIGGER IF EXISTS trigger__procedure_change_categories_order__swap_table_order;

DROP VIEW IF EXISTS procedure_remove_category;

DROP TRIGGER IF EXISTS trigger__procedure_remove_category__remove_category;

/*
 * Time periods
 */
DROP TABLE IF EXISTS time_periods;

DROP TRIGGER IF EXISTS dont_delete_from_time_periods;

DROP TRIGGER IF EXISTS dont_update_from_time_periods;

/*
 * Archived periods
 */
DROP TABLE IF EXISTS archived_periods;

/*
 * Input data to category and reset period
 */
DROP VIEW IF EXISTS procedure_fund_category_and_reset_period;

DROP TRIGGER IF EXISTS trigger__procedure_fund_category_and_reset_period_insert;

/*
 * Add money to category without reseting the period
 */
DROP VIEW IF EXISTS procedure_fund_category;

DROP TRIGGER IF EXISTS triggers__procedure_fund_category__insert;

/*
 * Enable foreign key constraints
 */
PRAGMA foreign_keys = ON;

/*
 * Set user version
 */
PRAGMA user_version = 0;