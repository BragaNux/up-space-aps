CREATE TABLE IF NOT EXISTS milestones (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('Motor', 'Linguagem', 'Social', 'Cognitivo')),
    description TEXT NOT NULL,
    achieved_at TIMESTAMPTZ NULL,
    done BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_milestones_student_id ON milestones(student_id);