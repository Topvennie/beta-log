-- name: UserGet :one
SELECT *
FROM users
WHERE id = $1;

-- name: UserGetByUID :one
SELECT *
FROM users
WHERE uid = $1;

-- name: UserGetAll :many
SELECT *
FROM users;

-- name: UserCreate :one
INSERT INTO users (uid, name)
VALUES ($1, $2)
RETURNING id;
