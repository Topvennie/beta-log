-- name: SessionGet :one
SELECT *
FROM sessions
WHERE id = $1;

-- name: SessionGetAll :many
SELECT *
FROM sessions
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: SessionCreate :one
INSERT INTO sessions (user_id, name, active, position)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: SessionUpdate :exec
UPDATE sessions
SET name = $2, active = $3, position = $4
WHERE id = $1;

-- name: SessionDelete :exec
UPDATE sessions
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
