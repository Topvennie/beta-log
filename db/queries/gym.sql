-- name: GymGet :one
SELECT *
FROM gyms
WHERE id = $1;

-- name: GymGetByExternalSource :one
SELECT *
FROM gyms
WHERE external_id = $1 AND source = $2;

-- name: GymGetAllByExternalSource :many
SELECT *
FROM gyms
WHERE external_id = ANY($1::int[]) AND source = $2;

-- name: GymGetAllByUser :many
SELECT *
FROM gyms
WHERE user_id = $1
ORDER BY name;

-- name: GymCreate :one
INSERT INTO gyms (user_id, external_id, name, icon_path, source)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GymUpdate :exec
UPDATE gyms
SET name = $2, icon_path = $3
WHERE id = $1;