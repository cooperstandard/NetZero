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

-- name: GetUnpaidDebtsByCreditorAndDebtor :many
SELECT
    debts.id
FROM
    debts
JOIN
    transactions
ON
    debts.transaction_id = transactions.id
WHERE
    debts.paid = FALSE AND $1 = debts.debtor AND $2 = debts.creditor AND transactions.group_id = $3;

-- name: PayDebts :one
UPDATE
    debts
SET
    paid = TRUE,
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    *;

-- name: PayDebtsByTransaction :many
UPDATE
    debts
SET
    paid = TRUE,
    updated_at = NOW()
WHERE
    transaction_id = $1
RETURNING
    *;

-- name: DeleteDebtById :one
DELETE FROM debts
    WHERE id = $1
RETURNING *;
