-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    mail CHAR(30),
    registered TIMESTAMP DEFAULT now() NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
