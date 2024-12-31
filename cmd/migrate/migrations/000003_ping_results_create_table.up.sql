-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Migration to create the `ping_results` table
CREATE TABLE IF NOT EXISTS ping_results (
    id TEXT PRIMARY KEY NOT NULL,
    monitor_id TEXT NOT NULL,
    status TEXT NOT NULL,
    response_time INTEGER NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (monitor_id) REFERENCES monitors (id) ON DELETE CASCADE
);
