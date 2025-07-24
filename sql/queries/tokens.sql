-- name: CreateToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, email, user_id, expires_at)
    VALUES ($1, NOW(), NOW(), $2, $3, $4)
RETURNING
    *;

-- name: GetToken :one
SELECT
    *
FROM
    refresh_tokens
WHERE
    token = $1
LIMIT 1;

-- name: RevokeToken :one
UPDATE
    refresh_tokens
SET
    revoked_at = $1,
    updated_at = NOW()
WHERE
    token = $2
RETURNING
    *;


