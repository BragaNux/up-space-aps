DROP TABLE IF EXISTS student_guardians;

ALTER TABLE students
    DROP COLUMN IF EXISTS guardian_user_id,
    DROP COLUMN IF EXISTS photo_url,
    DROP COLUMN IF EXISTS group_name,
    DROP COLUMN IF EXISTS teacher_name,
    DROP COLUMN IF EXISTS birth_date,
    DROP COLUMN IF EXISTS enrollment_code,
    DROP COLUMN IF EXISTS blood_type,
    DROP COLUMN IF EXISTS allergies,
    DROP COLUMN IF EXISTS restrictions,
    DROP COLUMN IF EXISTS medications;