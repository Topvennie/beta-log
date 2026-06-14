-- +goose Up
CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id),
  name TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT FALSE,
  position INTEGER,
  deleted_at TIMESTAMPTZ DEFAULT NULL,

  CONSTRAINT check_active_position CHECK (NOT active OR position IS NOT NULL)
);

-- +goose Down
DROP TABLE sessions;
