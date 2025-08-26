-- name: CreateBalance :one
INSERT INTO balances (user_id, group_id, creditor_id, updated_at, balance) 
  VALUES ($1, $2, $3, NOW(), $4)
RETURNING *;

-- name: GetBalanceForUserByGroup :many
SELECT balances.updated_at, balances.balance, users.name, users.email
  FROM balances JOIN users ON users.id = creditor.id
WHERE balances.group_id = $1 and user_id = $2;


