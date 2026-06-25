-- +goose Up
CREATE TABLE settings (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users (id),
  climb_toplogger_user_id TEXT,
  climb_toplogger_auth_token TEXT,
  climb_toplogger_refresh_token TEXT,
  climb_toplogger_expiration TIMESTAMPTZ
);

-- +goose Down
DROP TABLE settings;
