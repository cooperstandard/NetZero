-- name: CreateBalance :one
INSERT INTO balances (user_id, group_id, creditor_id, updated_at, balance) 
  VALUES ($1, $2, $3, NOW(), $4)
RETURNING *;

-- name: GetBalanceForDebtorByGroup :many
SELECT balances.updated_at, balances.balance, users.name, users.email
  FROM balances JOIN users ON users.id = creditor.id
WHERE balances.group_id = $1 and user_id = $2;

-- name: GetBalanceForCreditorByGroup :many
SELECT balances.updated_at, balances.balance, users.name, users.email
  FROM balances JOIN users ON users.id = creditor.id
WHERE balances.group_id = $1 and creditor_id = $2;

-- name: GetBalance :one
SELECT * FROM balances
WHERE
  user_id = $1 AND group_id = $2 and creditor_id = $3;

-- name: UpdateBalance :one
UPDATE
  balances
SET
  balance = $1, updated_at = NOW()
WHERE
  user_id = $2 AND group_id = $3 and creditor_id = $4
RETURNING *;

-- name: InsertOrUpdateBalance :one
INSERT INTO balances (user_id, group_id, creditor_id, updated_at, balance) VALUES ($1, $2, $3, NOW(), $4)
ON CONFLICT (user_id, group_id, creditor_id) DO UPDATE SET balance = balance + $4, updated_at = NOW()
RETURNING *;

