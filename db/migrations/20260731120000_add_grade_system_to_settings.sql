-- +goose Up
CREATE TYPE grade_system AS ENUM ('font', 'v');

ALTER TABLE settings ADD COLUMN grade_system GRADE_SYSTEM NOT NULL DEFAULT 'font';

-- +goose Down
ALTER TABLE settings DROP COLUMN grade_system;

DROP TYPE grade_system;
