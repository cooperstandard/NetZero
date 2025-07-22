-- name: CreateTransaction :one
INSERT INTO transactions (id, created_at, updated_at, title, description, author_id, amount)
    VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2, $3, $4)
RETURNING
    *;

