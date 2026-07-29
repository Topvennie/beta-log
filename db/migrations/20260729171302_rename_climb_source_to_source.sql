-- +goose Up
ALTER TYPE climb_source RENAME TO source;

-- +goose Down
ALTER TYPE source RENAME TO climb_source;
