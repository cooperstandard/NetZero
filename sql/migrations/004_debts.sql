-- +goose Up
CREATE TABLE debts (
    id uuid PRIMARY KEY,
    amount numeric(20, 2) NOT NULL,
    transaction_id uuid NOT NULL REFERENCES transactions (id) ON DELETE CASCADE,
    debtor uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    creditor uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    paid boolean NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE debts;

