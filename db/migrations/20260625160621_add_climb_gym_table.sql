-- +goose Up
CREATE TABLE climb_gyms (
  id SERIAL PRIMARY KEY,
  user_id  INT NOT NULL REFERENCES users (id),
  external_id TEXT NOT NULL,
  name TEXT NOT NULL,
  icon_path TEXT NOT NULL,

  UNIQUE (external_id)
);

-- +goose Down
DROP TABLE climb_gyms;
