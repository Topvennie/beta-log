-- name: ClimbGymGet :one
SELECT *
FROM climb_gyms
WHERE id = $1;

-- name: ClimbGymGetByExternalSource :one
SELECT *
FROM climb_gyms
WHERE external_id = $1 AND source = $2;

-- name: ClimbGymGetAllByExternalSource :many
SELECT *
FROM climb_gyms
WHERE external_id = ANY($1::int[]) AND source = $2;

-- name: ClimbGymCreate :one
INSERT INTO climb_gyms (user_id, external_id, name, icon_path, source)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: ClimbGymUpdate :exec
UPDATE climb_gyms
SET name = $2, icon_path = $3
WHERE id = $1;
