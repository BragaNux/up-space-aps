CREATE TABLE IF NOT EXISTS timeline_events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    icon TEXT NULL,
    category TEXT NULL,
    color TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE timeline_events ADD COLUMN IF NOT EXISTS icon TEXT NULL;
ALTER TABLE timeline_events ADD COLUMN IF NOT EXISTS category TEXT NULL;
ALTER TABLE timeline_events ADD COLUMN IF NOT EXISTS color TEXT NULL;

CREATE OR REPLACE VIEW timeline_items AS
SELECT
    id,
    title,
    to_char(occurred_at AT TIME ZONE 'America/Sao_Paulo', 'HH24:MI') AS "time",
    description,
    icon,
    category,
    color
FROM timeline_events;
