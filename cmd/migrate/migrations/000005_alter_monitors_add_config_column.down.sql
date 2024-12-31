-- SQLite does not support dropping columns directly.
-- To remove the `config` column, we need to create a new table without the column,
-- copy the data, drop the old table, and rename the new table.

-- Create a new table without the `config` column
CREATE TABLE IF NOT EXISTS monitors_new (
    id TEXT PRIMARY KEY NOT NULL,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    interval INTEGER NOT NULL,
    method TEXT,
    kind TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Copy data from the original table to the new table
INSERT INTO monitors_new (id, user_id, name, address, interval, method, kind, created_at, updated_at)
SELECT id, user_id, name, address, interval, method, kind, created_at, updated_at
FROM monitors;

-- Drop the original table
DROP TABLE monitors;

-- Rename the new table to the original name
ALTER TABLE monitors_new RENAME TO monitors;

-- Recreate the trigger on the updated table
CREATE TRIGGER IF NOT EXISTS update_monitors_updated_at
AFTER UPDATE ON monitors
FOR EACH ROW
BEGIN
    UPDATE monitors
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
