-- +goose Up
CREATE TABLE balances (
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
  creditor_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  updated_at timestamp NOT NULL,
  balance numeric(20, 2) NOT NULL,
  UNIQUE(user_id, group_id, creditor_id)
);

-- +goose Down
DROP TABLE balances;
