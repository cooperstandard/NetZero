-- +goose Up
CREATE TABLE balances (
  user_id UUID,
  group_id UUID,
  creditor_id UUID,
  updated_at timestamp NOT NULL,
  balance numeric(20, 2) NOT NULL,
  UNIQUE(user_id, group_id, creditor_id)
);


-- +goose Down
DROP TABLE balances;
