-- +goose Up
CREATE TABLE balances (
  user_id UUID NOT NULL,
  group_id UUID NOT NULL,
  creditor_id UUID NOT NULL,
  updated_at timestamp NOT NULL,
  balance numeric(20, 2) NOT NULL,
  UNIQUE(user_id, group_id, creditor_id)
);


-- +goose Down
DROP TABLE balances;
