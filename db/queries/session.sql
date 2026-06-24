-- name: SessionGet :many
SELECT sqlc.embed(s), sqlc.embed(se), sqlc.embed(e), sqlc.embed(v)
FROM sessions s
LEFT JOIN session_exercises_view se ON se.session_id = s.id
LEFT JOIN exercises_view e ON se.exercise_id = e.id
LEFT JOIN variants_view v ON se.variant_id = v.id
WHERE s.id = $1;

-- name: SessionGetAll :many
SELECT sqlc.embed(s), sqlc.embed(se), sqlc.embed(e), sqlc.embed(v)
FROM sessions s
LEFT JOIN session_exercises_view se ON se.session_id = s.id
LEFT JOIN exercises_view e ON se.exercise_id = e.id
LEFT JOIN variants_view v ON se.variant_id = v.id
WHERE s.user_id = $1 AND s.deleted_at IS NULL
ORDER BY s.name, s.id;

-- name: SessionGetByExercise :many
SELECT sqlc.embed(s), sqlc.embed(se), sqlc.embed(e), sqlc.embed(v)
FROM sessions s
LEFT JOIN session_exercises_view se ON se.session_id = s.id
LEFT JOIN exercises_view e ON se.exercise_id = e.id
LEFT JOIN variants_view v ON se.variant_id = v.id
WHERE e.id = $1;

-- name: SessionGetByVariants :many
SELECT sqlc.embed(s), sqlc.embed(se), sqlc.embed(e), sqlc.embed(v)
FROM sessions s
LEFT JOIN session_exercises_view se ON se.session_id = s.id
LEFT JOIN exercises_view e ON se.exercise_id = e.id
LEFT JOIN variants_view v ON se.variant_id = v.id
WHERE v.id = ANY($1::int[]);

-- name: SessionCreate :one
INSERT INTO sessions (user_id, name)
VALUES ($1, $2)
RETURNING id;

-- name: SessionUpdate :exec
UPDATE sessions
SET name = $2
WHERE id = $1;

-- name: SessionDelete :exec
UPDATE sessions
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
