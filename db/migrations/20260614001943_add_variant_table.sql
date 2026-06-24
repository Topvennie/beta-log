-- +goose Up
CREATE TABLE variants (
  id SERIAL PRIMARY KEY,
  exercise_id INTEGER NOT NULL REFERENCES exercises (id),
  variant TEXT NOT NULL
);

CREATE VIEW variants_view AS SELECT variants.* FROM (SELECT NULL WHERE false) dummy FULL JOIN variants ON true;

-- +goose Down
DROP VIEW variants_view;

DROP TABLE variants;
