-- +goose Up
ALTER TABLE climbs
DROP COLUMN color;

-- +goose Down
ALTER TABLE climbs
ADD COLUMN color TEXT;
