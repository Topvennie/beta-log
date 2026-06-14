-- +goose Up
CREATE TABLE sessions_exercises (
  id SERIAL PRIMARY KEY,
  session_id INTEGER NOT NULL REFERENCES sessions (id),
  exercise_id INTEGER NOT NULL REFERENCES exercises (id),
  position INTEGER NOT NULL,
  sets INTEGER NOT NULL,
  reps INTEGER,
  weight INTEGER,
  duration_s INTEGER,

  UNIQUE (session_id, exercise_id)
);

-- +goose Down
DROP TABLE sessions_exercises;
