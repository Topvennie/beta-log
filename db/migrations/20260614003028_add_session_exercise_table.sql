-- +goose Up
CREATE TABLE session_exercises (
  id SERIAL PRIMARY KEY,
  session_id INTEGER NOT NULL REFERENCES sessions (id),
  exercise_id INTEGER NOT NULL REFERENCES exercises (id),
  variant_id INTEGER REFERENCES variants (id),
  position INTEGER NOT NULL,
  sets INTEGER NOT NULL,
  reps INTEGER,
  weight INTEGER,
  duration_s INTEGER,

  UNIQUE (session_id, exercise_id)
);

CREATE VIEW session_exercises_view AS SELECT session_exercises.* FROM (SELECT NULL WHERE false) dummy FULL JOIN session_exercises ON true;

-- +goose Down
DROP VIEW session_exercises_view;

DROP TABLE session_exercises;
