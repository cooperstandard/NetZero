-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    name TEXT,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    email text NOT NULL UNIQUE,
    hashed_password text NOT NULL
);

-- +goose Down
DROP TABLE users;

