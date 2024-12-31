-- Create the status_pages table
CREATE TABLE IF NOT EXISTS status_pages (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to automatically update `updated_at` timestamp on record update
CREATE TRIGGER IF NOT EXISTS update_status_pages_updated_at
AFTER UPDATE ON status_pages
FOR EACH ROW
BEGIN
    UPDATE status_pages
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;

-- Create a join table to handle the many-to-many relationship between status_pages and monitors
CREATE TABLE IF NOT EXISTS status_page_monitors (
    status_page_id TEXT NOT NULL,
    monitor_id TEXT NOT NULL,
    PRIMARY KEY (status_page_id, monitor_id),
    FOREIGN KEY (status_page_id) REFERENCES status_pages (id) ON DELETE CASCADE,
    FOREIGN KEY (monitor_id) REFERENCES monitors (id) ON DELETE CASCADE
);
