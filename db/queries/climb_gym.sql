-- name: ClimbGymGet :one
SELECT *
FROM climb_gyms
WHERE id = $1;

-- name: ClimbGymGetByExternal :one
SELECT *
FROM climb_gyms
WHERE external_id = $1;

-- name: ClimbGymGetByExternalIds :many
SELECT *
FROM climb_gyms
WHERE external_id = ANY($1::int[]);

-- name: ClimbGymCreate :one
INSERT INTO climb_gyms (user_id, external_id, name, icon_path)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: ClimbGymUpdate :exec
UPDATE climb_gyms
SET name = $2, icon_path = $3
WHERE id = $1;
