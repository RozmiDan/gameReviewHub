-- +goose Up
CREATE TABLE IF NOT EXISTS games(
    id SERIAL PRIMARY KEY,
    gamedescription TEXT NOT NULL UNIQUE,
    gamename CHAR(30),
    registered TIMESTAMP DEFAULT now() NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS games;
