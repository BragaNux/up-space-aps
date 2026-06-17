CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    avatar TEXT NULL,
    classroom TEXT NULL,
    presence_status TEXT NOT NULL DEFAULT 'absent',
    check_in_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE students ADD COLUMN IF NOT EXISTS avatar TEXT NULL;
ALTER TABLE students ADD COLUMN IF NOT EXISTS classroom TEXT NULL;
