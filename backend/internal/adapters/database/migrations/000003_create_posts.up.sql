CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    author_name TEXT NOT NULL DEFAULT 'Equipe UP Espaco',
    author_avatar TEXT NULL,
    author_role TEXT NOT NULL DEFAULT 'Educador(a)',
    title TEXT NOT NULL,
    category TEXT NULL,
    description TEXT NOT NULL,
    image_url TEXT NULL,
    pedagogical_note TEXT NOT NULL,
    likes BIGINT NOT NULL DEFAULT 0,
    comments_count BIGINT NOT NULL DEFAULT 0,
    bookmarks BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE posts ADD COLUMN IF NOT EXISTS author_name TEXT NOT NULL DEFAULT 'Equipe UP Espaco';
ALTER TABLE posts ADD COLUMN IF NOT EXISTS author_avatar TEXT NULL;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS author_role TEXT NOT NULL DEFAULT 'Educador(a)';
ALTER TABLE posts ADD COLUMN IF NOT EXISTS category TEXT NULL;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS image_url TEXT NULL;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS comments_count BIGINT NOT NULL DEFAULT 0;

CREATE OR REPLACE VIEW activity_posts AS
SELECT
    id,
    author_name,
    author_avatar,
    author_role,
    to_char(created_at AT TIME ZONE 'America/Sao_Paulo', 'HH24:MI') AS "time",
    title,
    category,
    description,
    image_url AS "imageUrl",
    likes,
    comments_count AS "commentsCount"
FROM posts;
