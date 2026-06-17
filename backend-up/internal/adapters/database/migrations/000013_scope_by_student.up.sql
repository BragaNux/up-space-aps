ALTER TABLE posts ADD COLUMN IF NOT EXISTS student_id INTEGER REFERENCES students(id) ON DELETE CASCADE;
UPDATE posts SET student_id = (SELECT id FROM students ORDER BY id ASC LIMIT 1) WHERE student_id IS NULL;
ALTER TABLE posts ALTER COLUMN student_id SET NOT NULL;

ALTER TABLE timeline_events ADD COLUMN IF NOT EXISTS student_id INTEGER REFERENCES students(id) ON DELETE CASCADE;
UPDATE timeline_events SET student_id = (SELECT id FROM students ORDER BY id ASC LIMIT 1) WHERE student_id IS NULL;
ALTER TABLE timeline_events ALTER COLUMN student_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_posts_student_id ON posts(student_id);
CREATE INDEX IF NOT EXISTS idx_timeline_events_student_id ON timeline_events(student_id);
