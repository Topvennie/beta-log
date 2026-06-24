-- name: SessionExerciseGetBySession :many
SELECT *
FROM session_exercises
WHERE session_id = $1
ORDER BY position;

-- name: SessionExerciseCreate :one
INSERT INTO session_exercises (session_id, exercise_id, variant_id, position, sets, reps, weight, duration_s)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: SessionExerciseDelete :exec
DELETE FROM session_exercises
WHERE id = $1;

-- name: SessionExerciseDeleteBySession :exec
DELETE FROM session_exercises
WHERE session_id = $1;
