-- +goose Up
ALTER TABLE climb_gyms RENAME TO gyms;

ALTER TABLE gyms RENAME CONSTRAINT climb_gyms_external_id_source_key TO gyms_external_id_source_key;

-- +goose Down
ALTER TABLE gyms RENAME CONSTRAINT gyms_external_id_source_key TO climb_gyms_external_id_source_key;

ALTER TABLE gyms RENAME TO climb_gyms;