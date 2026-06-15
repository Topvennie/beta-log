-- +goose Up
CREATE TABLE exercises (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id),
  name TEXT NOT NULL,
  variants TEXT[] NOT NULL,
  deleted_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE exercises;
