-- +goose Up
CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id),
  name TEXT NOT NULL,
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- +goose Down
DROP TABLE sessions;
