-- +goose Up
CREATE TYPE climb_source AS ENUM ('toplogger');

ALTER TABLE climb_gyms ADD COLUMN source CLIMB_SOURCE NOT NULL DEFAULT 'toplogger';
ALTER TABLE climb_days ADD COLUMN source CLIMB_SOURCE NOT NULL DEFAULT 'toplogger';
ALTER TABLE climbs ADD COLUMN source CLIMB_SOURCE NOT NULL DEFAULT 'toplogger';

ALTER TABLE climb_gyms DROP CONSTRAINT climb_gyms_external_id_key;
ALTER TABLE climb_days DROP CONSTRAINT climb_days_external_id_key;

ALTER TABLE climb_gyms ADD CONSTRAINT climb_gyms_external_id_source_key UNIQUE (external_id, source);
ALTER TABLE climb_days ADD CONSTRAINT climb_days_external_id_source_key UNIQUE (external_id, source);

ALTER TABLE climb_gyms ALTER COLUMN source DROP DEFAULT;
ALTER TABLE climb_days ALTER COLUMN source DROP DEFAULT;
ALTER TABLE climbs ALTER COLUMN source DROP DEFAULT;

-- +goose Down
ALTER TABLE climbs ALTER COLUMN source SET DEFAULT 'toplogger';
ALTER TABLE climb_days ALTER COLUMN source SET DEFAULT 'toplogger';
ALTER TABLE climb_gyms ALTER COLUMN source SET DEFAULT 'toplogger';

ALTER TABLE climb_days DROP CONSTRAINT climb_days_external_id_source_key;
ALTER TABLE climb_gyms DROP CONSTRAINT climb_gyms_external_id_source_key;

ALTER TABLE climb_days ADD CONSTRAINT climb_days_external_id_key UNIQUE (external_id);
ALTER TABLE climb_gyms ADD CONSTRAINT climb_gyms_external_id_key UNIQUE (external_id);

ALTER TABLE climbs DROP COLUMN source;
ALTER TABLE climb_days DROP COLUMN source;
ALTER TABLE climb_gyms DROP COLUMN source;

DROP TYPE climb_source;
