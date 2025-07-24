-- +goose up
CREATE TABLE refresh_tokens (
    token text PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    email text NOT NULL,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at timestamp NOT NULL,
    revoked_at timestamp
);

-- +goose down
DROP TABLE refresh_tokens;

