-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Migration to create the `monitors` table
CREATE TABLE IF NOT EXISTS monitors (
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

-- Trigger to automatically update `updated_at` timestamp on record update
CREATE TRIGGER IF NOT EXISTS update_monitors_updated_at
AFTER UPDATE ON monitors
FOR EACH ROW
BEGIN
    UPDATE monitors
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
