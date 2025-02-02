-- Revert the change by renaming the "password" column back to "password_hash".
ALTER TABLE users RENAME COLUMN password TO password_hash;
