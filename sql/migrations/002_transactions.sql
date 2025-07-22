
-- +goose Up
CREATE TABLE transactions (
  id        UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title      TEXT  NOT NULL,
  description TEXT,
  author_id UUID REFERENCES users(id) ON DELETE CASCADE,
  -- group_id TEXT NOT NULL, -- TODO: groups and debts
  amount NUMERIC(11,2) NOT NULL
  
);

-- +goose Down
DROP TABLE transactions;
