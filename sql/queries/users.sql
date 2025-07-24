-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, name)
    VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2, $3)
RETURNING
    *;

-- name: GetUsers :many
SELECT
    *
FROM
    users;

-- name: RemoveAllUsers :exec
DELETE FROM users;

