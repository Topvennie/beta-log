-- name: SessionExerciseGet :one
SELECT *
FROM sessions_exercises
WHERE id = $1;

-- name: SessionExerciseGetBySession :many
SELECT *
FROM sessions_exercises
WHERE session_id = $1
ORDER BY position;

-- name: SessionExerciseCreate :one
INSERT INTO sessions_exercises (session_id, exercise_id, variant, position, sets, reps, weight, duration_s)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: SessionExerciseUpdate :exec
UPDATE sessions_exercises
SET variant = $2, position = $3, sets = $4, reps = $5, weight = $6, duration_s = $7
WHERE id = $1;

-- name: SessionExerciseDelete :exec
DELETE FROM sessions_exercises
WHERE id = $1;

-- name: SessionExerciseDeleteBySession :exec
DELETE FROM sessions_exercises
WHERE session_id = $1;
