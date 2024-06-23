/*
 * Disable foreign key constraints just in case
 */
PRAGMA foreign_keys = OFF;

/*
 * User status table
 */
DROP TRIGGER IF EXISTS dont_delete_from_user_status;

DROP TRIGGER IF EXISTS dont_update_from_user_status;

DROP TABLE IF EXISTS user_status;

/*
 * Users table
 *
 * Make db_node field Nullable
 * Add status column
 */
-- disable foreign key constraint check
PRAGMA foreign_keys = off;

-- Here you can drop column
CREATE TABLE
    IF NOT EXISTS users_copy (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        user_email TEXT UNIQUE NOT NULL,
        db_version INTEGER,
        db_node INTEGER REFERENCES db_nodes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT null
    );

-- copy data from the table to the new_table
INSERT INTO
    users_copy (
        id,
        user_email,
        db_version,
        db_node,
        created_at,
        updated_at
    )
SELECT
    id,
    user_email,
    db_version,
    db_node,
    created_at,
    updated_at
FROM
    users;

-- drop the table
DROP TABLE users;

-- rename the new_table to the table
ALTER TABLE users_copy
RENAME TO users;

/*
 * Enable foreign key constraints
 */
PRAGMA foreign_keys = ON;

/*
 * Set user version
 */
PRAGMA user_version = 1;