-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
  id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  game_id    UUID        NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  user_id    UUID        NOT NULL,
  text       TEXT        NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_comments_game_id 
  ON comments(game_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS comments;
