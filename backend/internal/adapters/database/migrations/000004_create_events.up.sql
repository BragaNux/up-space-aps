CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    location TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    is_important BOOLEAN NOT NULL DEFAULT false,
    icon TEXT NULL,
    color TEXT NULL,
    rsvp_prompt TEXT NULL,
    is_rsvped BOOLEAN NOT NULL DEFAULT false,
    rsvp_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE events ADD COLUMN IF NOT EXISTS is_important BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE events ADD COLUMN IF NOT EXISTS icon TEXT NULL;
ALTER TABLE events ADD COLUMN IF NOT EXISTS color TEXT NULL;
ALTER TABLE events ADD COLUMN IF NOT EXISTS rsvp_prompt TEXT NULL;
ALTER TABLE events ADD COLUMN IF NOT EXISTS is_rsvped BOOLEAN NOT NULL DEFAULT false;

CREATE OR REPLACE VIEW calendar_events AS
SELECT
    id,
    title,
    to_char(starts_at AT TIME ZONE 'America/Sao_Paulo', 'DD/MM HH24:MI') AS "time",
    is_important AS "isImportant",
    icon,
    color,
    rsvp_prompt AS "rsvpPrompt",
    is_rsvped AS "isRSVPed"
FROM events;
