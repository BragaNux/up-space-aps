CREATE TABLE IF NOT EXISTS announcements (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    sender TEXT NOT NULL,
    priority TEXT NOT NULL CHECK (priority IN ('Urgente', 'Importante', 'Informativo')),
    preview TEXT NOT NULL,
    body TEXT NOT NULL,
    attachment_name TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS announcement_reads (
    id SERIAL PRIMARY KEY,
    announcement_id INTEGER NOT NULL REFERENCES announcements(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    read_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (announcement_id, user_id)
);