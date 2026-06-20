-- name: ExerciseGet :one
SELECT *
FROM exercises
WHERE id = $1;

-- name: ExerciseGetAll :many
SELECT *
FROM exercises
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY name, variants;

-- name: ExerciseGetByIDs :many
SELECT *
FROM exercises
WHERE id = ANY($1::int[]) AND deleted_at IS NULL;

-- name: ExerciseCreate :one
INSERT INTO exercises (user_id, name, variants)
VALUES ($1, $2, $3)
RETURNING id;

-- name: ExerciseUpdate :exec
UPDATE exercises
SET name = $2, variants = $3
WHERE id = $1;

-- name: ExerciseDelete :exec
UPDATE exercises
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
