-- name: CreateGroup :one
INSERT INTO GROUPS (id, created_at, updated_at, name)
    VALUES (gen_random_uuid (), NOW(), NOW(), $1)
RETURNING
    *;

-- name: JoinGroup :one
INSERT INTO group_members (user_id, group_id)
    VALUES ($1, $2)
RETURNING
    *;

-- name: GetGroupsByUser :many
SELECT
    id,
    name,
    created_at
FROM
    GROUPS
    JOIN group_members ON groups.id = group_members.group_id
WHERE
    user_id = $1;

-- name: GetUsersByGroup :many
SELECT
    id,
    name
FROM
    users
    JOIN group_members ON users.id = group_members.user_id
WHERE
    group_id = $1
    AND id != $2;

-- name: GetGroupByName :one
SELECT * FROM groups WHERE name = $1;
