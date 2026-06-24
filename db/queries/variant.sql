-- name: VariantGet :one
SELECT *
FROM variants
WHERE id = $1;

-- name: VariantGetByExercise :many
SELECT *
FROM variants
WHERE exercise_id = $1;

-- name: VariantGetByIDs :many
SELECT *
FROM variants
WHERE id = ANY($1::int[]);

-- name: VariantCreate :one
INSERT INTO variants (exercise_id, variant)
VALUES ($1, $2)
RETURNING id;

-- name: VariantUpdate :exec
UPDATE variants
SET variant = $2
WHERE id = $1;

-- name: VariantDelete :exec
DELETE FROM variants
WHERE id = $1;

-- name: VariantDeleteByExercise :exec
DELETE FROM variants
WHERE exercise_id = $1;

