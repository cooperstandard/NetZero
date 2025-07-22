-- +goose Up
CREATE TABLE debts (
  id UUID PRIMARY KEY,
  amount NUMERIC(11,2) NOT NULL,
  transaction_id UUID NOT NULL REFERENCES transactions(id),
  debtor UUID NOT NULL REFERENCES users(id),
  creditor UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  paid BOOLEAN NOT NULL DEFAULT false


);

-- +goose Down
DROP TABLE debts;
