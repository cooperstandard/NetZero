-- +goose up
ALTER TABLE users ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';
