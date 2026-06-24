-- name: ExerciseGet :many
SELECT sqlc.embed(e), sqlc.embed(ev)
FROM exercises e
LEFT JOIN variants_view ev ON e.id = ev.exercise_id
WHERE e.id = $1;

-- name: ExerciseGetAll :many
SELECT sqlc.embed(e), sqlc.embed(ev)
FROM exercises e
LEFT JOIN variants_view ev ON e.id = ev.exercise_id
WHERE e.user_id = $1 AND e.deleted_at IS NULL
ORDER BY e.name, ev.id;

-- name: ExerciseGetByIDs :many
SELECT sqlc.embed(e), sqlc.embed(ev)
FROM exercises e
LEFT JOIN variants_view ev ON e.id = ev.exercise_id
WHERE e.id = ANY($1::int[]) AND e.deleted_at IS NULL
ORDER BY e.name, ev.id;

-- name: ExerciseCreate :one
INSERT INTO exercises (user_id, name)
VALUES ($1, $2)
RETURNING id;

-- name: ExerciseUpdate :exec
UPDATE exercises
SET name = $2
WHERE id = $1;

-- name: ExerciseDelete :exec
UPDATE exercises
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
