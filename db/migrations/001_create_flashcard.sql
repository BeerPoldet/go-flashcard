-- +goose Up
CREATE TABLE IF NOT EXISTS flashcards (
  id SERIAL PRIMARY KEY,
  title text NOT NULL,
  content text NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE flashcards;
