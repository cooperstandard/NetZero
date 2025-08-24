-- +goose Up
CREATE TABLE transactions (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title text NOT NULL,
    description text,
    author_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    group_id uuid NOT NULL REFERENCES GROUPS (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE transactions;

