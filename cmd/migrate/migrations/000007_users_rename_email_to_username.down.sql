-- Revert the change by renaming the "username" column back to "email".
ALTER TABLE users RENAME COLUMN username TO email;
