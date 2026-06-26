-- name: ClimbDayGet :one
SELECT *
FROM climb_days
WHERE id = $1;

-- name: ClimbDayGetByExternal :one
SELECT *
FROM climb_days
WHERE external_id = $1;

-- name: ClimbDayGetPopulated :many
SELECT sqlc.embed(d), sqlc.embed(c), sqlc.embed(g)
FROM climb_days d
LEFT  JOIN climbs c ON c.climb_day_id = d.id
LEFT JOIN climb_gyms g ON d.gym_id = g.id
WHERE d.id = $1;

-- name: ClimbDayGetPopulatedByExternal :many
SELECT sqlc.embed(d), sqlc.embed(c), sqlc.embed(g)
FROM climb_days d
LEFT  JOIN climbs c ON c.climb_day_id = d.id
LEFT JOIN climb_gyms g ON d.gym_id = g.id
WHERE d.external_id = $1;

-- name: ClimbDayGetAllPopulatedByExternal :many
SELECT sqlc.embed(d), sqlc.embed(c), sqlc.embed(g)
FROM climb_days d
LEFT  JOIN climbs c ON c.climb_day_id = d.id
LEFT JOIN climb_gyms g ON d.gym_id = g.id
WHERE d.external_id = ANY($1::int[]);

-- name: ClimbDayCreate :one
INSERT INTO climb_days (user_id, external_id, gym_id, date)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: ClimbDayUpdate :exec
UPDATE climb_days
SET gym_id = $2, date = $3
WHERE id = $1;
