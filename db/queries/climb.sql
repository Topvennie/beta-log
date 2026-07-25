-- name: ClimbGet :one
SELECT *
FROM climbs
WHERE id = $1;

-- name: ClimbGetByExternal :one
SELECT *
FROM climbs
WHERE external_id = $1;

-- name: ClimbGetAllByClimbDay :many
SELECT *
FROM climbs
WHERE climb_day_id = $1;

-- name: ClimbCreate :one
INSERT INTO climbs (user_id, external_id, climb_day_id, grade, hold_color, climb_type, finish_type, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: ClimbUpdate :exec
UPDATE climbs
SET grade = $2, hold_color = $3, climb_type = $4, finish_type = $5
WHERE id = $1;
