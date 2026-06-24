-- +goose Up
CREATE TABLE exercises (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id),
  name TEXT NOT NULL,
  deleted_at TIMESTAMPTZ
);

CREATE VIEW exercises_view AS SELECT exercises.* FROM (SELECT NULL WHERE false) dummy FULL JOIN exercises ON true;

-- +goose Down
DROP VIEW exercises_view;

DROP TABLE exercises;
