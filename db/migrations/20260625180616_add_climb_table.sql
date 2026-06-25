-- +goose Up
CREATE TYPE climb_type AS ENUM ('boulder', 'lead');
CREATE TYPE finish_type AS ENUM ('flash', 'top', 'repeat');

CREATE TABLE climbs (
  id SERIAL PRIMARY KEY,
  user_id  INT NOT NULL REFERENCES users (id),
  external_id TEXT NOT NULL,
  climb_day_id INT NOT NULL REFERENCES climb_days (id),
  grade INT NOT NULL,
  color TEXT NOT NULL,
  hold_color TEXT NOT NULL,
  climb_type CLIMB_TYPE NOT NULL,
  finish_type FINISH_TYPE NOT NULL
);

-- +goose Down
DROP TABLE climbs;

DROP TYPE finish_type;
DROP TYPE climb_type;
