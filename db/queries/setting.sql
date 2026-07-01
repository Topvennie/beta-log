-- name: SettingGetByUser :one
SELECT *
FROM settings
WHERE user_id = $1;

-- name: SettingCreate :one
INSERT INTO settings (user_id)
VALUES ($1)
RETURNING id;

-- name: SettingToploggerUpdate :exec
UPDATE settings
SET climb_toplogger_user_id = $2, climb_toplogger_auth_token = $3, climb_toplogger_refresh_token = $4, climb_toplogger_expiration = $5
WHERE id = $1;
