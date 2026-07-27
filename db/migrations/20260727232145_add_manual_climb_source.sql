-- +goose Up
ALTER TYPE climb_source ADD VALUE 'manual';

-- +goose Down
-- PostgreSQL does not support removing a value from an enum.
