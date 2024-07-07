/*
 * Users table
 *
 * Add password_hash column
 */
-- disable foreign key constraint check
PRAGMA foreign_keys = off;

-- Here you can drop column
CREATE TABLE
    IF NOT EXISTS users_copy (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        user_email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        db_version INTEGER,
        db_node INTEGER REFERENCES db_nodes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        status INTEGER NOT NULL REFERENCES user_status ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT null
    );

-- copy data from the table to the new_table
INSERT INTO
    users_copy (
        id,
        user_email,
        password_hash,
        db_version,
        db_node,
        status,
        created_at,
        updated_at
    )
SELECT
    id,
    user_email,
    "",
    db_version,
    db_node,
    status,
    created_at,
    updated_at
FROM
    users;

-- drop the table
DROP TABLE users;

-- rename the new_table to the table
ALTER TABLE users_copy
RENAME TO users;

-- enable foreign key constraint check
PRAGMA foreign_keys = on;

/*
 * Set user version
 */
PRAGMA user_version = 3;