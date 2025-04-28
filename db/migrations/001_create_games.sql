-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS games (
  id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  name         TEXT        NOT NULL,
  genre        TEXT        NOT NULL DEFAULT 'Unknown',
  creator      TEXT        NOT NULL DEFAULT 'Unknown',
  description  TEXT        NOT NULL DEFAULT '',
  release_date DATE        NOT NULL,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);


-- +goose Down
DROP TABLE IF EXISTS games;
