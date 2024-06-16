/*
 * Users table
 *
 * Contains a list of all users
 * with email and db node reference that contains the sqlite db, along with db version
 */
CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        user_email TEXT UNIQUE NOT NULL,
        db_version INTEGER,
        db_node INTEGER NOT NULL REFERENCES db_nodes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT null
    );

/*
 * DB nodes table
 *
 * Contains a list of all running db nodes
 */
CREATE TABLE
    IF NOT EXISTS db_nodes (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        remote_address TEXT UNIQUE NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT null
    );

/*
 * Set user version to 1
 */
PRAGMA user_version = 1;