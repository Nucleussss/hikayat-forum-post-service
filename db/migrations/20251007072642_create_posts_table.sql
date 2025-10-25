-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  author_id UUID NOT NULL,           
  category VARCHAR(50) NOT NULL DEFAULT 'general',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_category ON posts(category);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_posts_author_id;
DROP INDEX idx_posts_category;
DROP INDEX idx_posts_created_at;
DROP TABLE posts;
-- +goose StatementEnd
