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
INSERT INTO sessions_exercises (session_id, exercise_id, position, sets, reps, weight, duration_s)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: SessionExerciseUpdate :exec
UPDATE sessions_exercises
SET position = $2, sets = $3, reps = $4, weight = $5, duration_s = $6
WHERE id = $1;

-- name: SessionExerciseDelete :exec
DELETE FROM sessions_exercises
WHERE id = $1;
