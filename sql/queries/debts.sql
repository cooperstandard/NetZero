-- name: CreateDebt :one
INSERT INTO debts (id, amount, transaction_id, debtor, creditor, created_at, updated_at)
    VALUES (gen_new_uuid(), $1, $2, $3, $4, NOW(), NOW())
RETURNING
    *;

-- name: GetDebtsByTransaction :many
SELECT
    *
FROM
    debts
WHERE
    $1 = transaction_id;

-- name: GetDebtsByDebtor :many
SELECT
    *
FROM
    debts
WHERE
    $1 = debtor;

-- name: GetDebtsByCreditor :many
SELECT
    *
FROM
    debts
WHERE
    $1 = creditor;

-- name: PayDebts :one
UPDATE
    debts
SET
    paid = TRUE
WHERE
    id = $1
RETURNING
    *;

-- name: PayDebtsByTransaction :many
UPDATE
    debts
SET
    paid = TRUE
WHERE
    transaction_id = $1
RETURNING
    *;

-- name: DeleteDebtById :one
DELETE FROM debts
    WHERE id = $1
RETURNING *;
