-- name: ExerciseGet :one
SELECT *
FROM exercises
WHERE id = $1;

-- name: ExerciseGetAll :many
SELECT *
FROM exercises
WHERE user_id = $1 AND NOT DELETED
ORDER BY name;

-- name: ExerciseGetByIDs :many
SELECT *
FROM exercises
WHERE id = ANY($1::int[]) AND NOT DELETED;

-- name: ExerciseCreate :one
INSERT INTO exercises (user_id, name, variant)
VALUES ($1, $2, $3)
RETURNING id;

-- name: ExerciseUpdate :exec
UPDATE exercises
SET name = $2, variant = $3
WHERE id = $1;

-- name: ExerciseDelete :exec
UPDATE exercises
SET deleted_at = NOW()
WHERE id = $1 AND NOT DELETED;
