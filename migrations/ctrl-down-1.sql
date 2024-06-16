/*
 * Disable foreign key constraints just in case
 */
PRAGMA foreign_keys = OFF;

/*
 * Remove users table
 */
DROP TABLE IF EXISTS users;

/*
 * Remove DB nodes table
 */
DROP TABLE IF EXISTS db_nodes;

/*
 * Set user version to 0
 */
PRAGMA user_version = 0;

/*
 * Enable foreign key constraints
 */
PRAGMA foreign_keys = ON;