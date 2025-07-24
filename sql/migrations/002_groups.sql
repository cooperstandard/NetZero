-- +goose Up
CREATE TABLE GROUPS (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name text NOT NULL
);

-- +goose Down
DROP TABLE GROUPS;

