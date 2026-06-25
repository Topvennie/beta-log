-- +goose Up
CREATE TABLE climb_days (
  id SERIAL PRIMARY KEY,
  user_id  INT NOT NULL REFERENCES users (id),
  external_id TEXT NOT NULL,
  gym_id INT NOT NULL REFERENCES climb_gyms (id),
  date timestamptz NOT NULL,

  UNIQUE (external_id)
);

-- +goose Down
DROP TABLE climb_days;
