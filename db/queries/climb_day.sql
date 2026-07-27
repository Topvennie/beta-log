-- name: ClimbDayGet :one
SELECT *
FROM climb_days
WHERE id = $1;

-- name: ClimbDayGetByExternalSource :one
SELECT *
FROM climb_days
WHERE external_id = $1 AND source = $2;

-- name: ClimbDayGetPopulated :many
SELECT sqlc.embed(d), sqlc.embed(c), sqlc.embed(g)
FROM climb_days d
LEFT  JOIN climbs c ON c.climb_day_id = d.id
LEFT JOIN gyms g ON d.gym_id = g.id
WHERE d.id = $1;

-- name: ClimbDayGetPopulatedByExternalSource :many
SELECT sqlc.embed(d), sqlc.embed(c), sqlc.embed(g)
FROM climb_days d
LEFT  JOIN climbs c ON c.climb_day_id = d.id
LEFT JOIN gyms g ON d.gym_id = g.id
WHERE d.external_id = $1 AND d.source = $2;

-- name: ClimbDayGetAllPopulatedFiltered :many
SELECT sqlc.embed(d), sqlc.embed(c), sqlc.embed(g)
FROM climb_days d
LEFT JOIN climbs c ON c.climb_day_id = d.id
LEFT JOIN gyms g ON d.gym_id = g.id
WHERE d.id IN (
  SELECT id FROM climb_days AS cd
  WHERE cd.user_id = $1
  ORDER BY cd.date DESC
  LIMIT $2 OFFSET $3
)
ORDER BY d.date DESC, c.grade DESC;

-- name: ClimbDayGetAllPopulatedByExternalSource :many
SELECT sqlc.embed(d), sqlc.embed(c), sqlc.embed(g)
FROM climb_days d
LEFT  JOIN climbs c ON c.climb_day_id = d.id
LEFT JOIN gyms g ON d.gym_id = g.id
WHERE d.external_id = ANY($1::int[]) AND d.source = $2
ORDER BY d.date ASC;

-- name: ClimbDayGetAllPopulatedByUser :many
SELECT sqlc.embed(d), sqlc.embed(c), sqlc.embed(g)
FROM climb_days d
LEFT  JOIN climbs c ON c.climb_day_id = d.id
LEFT JOIN gyms g ON d.gym_id = g.id
WHERE d.user_id = $1
ORDER BY d.date ASC;

-- name: ClimbDayCreate :one
INSERT INTO climb_days (user_id, external_id, gym_id, date, source)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: ClimbDayUpdate :exec
UPDATE climb_days
SET gym_id = $2, date = $3
WHERE id = $1;
