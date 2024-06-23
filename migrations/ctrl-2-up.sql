/*
 * User status table
 *
 * new - when a user is created and not assigned to a DB Node
 * assigned - when a user is assigned to a db node
 */
CREATE TABLE
    IF NOT EXISTS user_status (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE CHECK (name IN ('new', 'assigned'))
    );

/*
 * Don't allow deletion from user_status
 */
CREATE TRIGGER dont_delete_from_user_status BEFORE DELETE ON user_status BEGIN
SELECT
    RAISE (ABORT, 'Cant delete from user_status');

END;

/*
 * Don't allow updates in user_status
 */
CREATE TRIGGER dont_update_from_user_status BEFORE
UPDATE ON user_status BEGIN
SELECT
    RAISE (ABORT, 'Cant update in user_status');

END;

/*
 * Insert values
 */
INSERT INTO
    user_status (name)
VALUES
    ('new'),
    ('assigned');

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
        status INTEGER NOT NULL REFERENCES user_status ON UPDATE CASCADE ON DELETE RESTRICT,
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
        status,
        created_at,
        updated_at
    )
SELECT
    id,
    user_email,
    db_version,
    db_node,
    (
        CASE
            WHEN db_node != NULL THEN 2
            ELSE 1
        END
    ) as 'status',
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
PRAGMA user_version = 2;