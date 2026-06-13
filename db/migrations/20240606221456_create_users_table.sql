-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
