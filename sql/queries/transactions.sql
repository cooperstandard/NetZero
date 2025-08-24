-- name: CreateTransaction :one
INSERT INTO transactions (id, created_at, updated_at, title, description, author_id, group_id)
    VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2, $3, $4)
RETURNING
    *;

-- name: GetTransactionsByGroup :many
SELECT
    *
FROM
    transactions
WHERE
    group_id = $1;

-- name: GetTransactonsByAuthor :many
SELECT
    *
FROM
    transactions
WHERE
    author_id = $1;

