-- +goose Up
CREATE TYPE task_result AS ENUM ('success', 'failed');

CREATE TABLE tasks (
  uid VARCHAR(255) NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  active BOOLEAN NOT NULL,
  recurring BOOLEAN NOT NULL
);

CREATE TABLE task_runs (
  id SERIAL PRIMARY KEY,
  task_uid VARCHAR(255) NOT NULL REFERENCES tasks (uid) ON DELETE CASCADE,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  run_at TIMESTAMPTZ NOT NULL,
  result TASK_RESULT NOT NULL,
  error TEXT,
  duration BIGINT NOT NULL,
  message TEXT
);

-- +goose Down
DROP TABLE task_runs;
DROP TABLE tasks;

DROP TYPE task_result;
