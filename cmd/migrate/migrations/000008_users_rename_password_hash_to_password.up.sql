-- Rename the "password_hash" column to "password" in the "users" table.
ALTER TABLE users RENAME COLUMN password_hash TO password;
