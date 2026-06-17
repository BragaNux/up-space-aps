CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id INTEGER NULL REFERENCES users(id) ON DELETE SET NULL,
    author_name TEXT NOT NULL,
    avatar_url TEXT NOT NULL DEFAULT '',
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);