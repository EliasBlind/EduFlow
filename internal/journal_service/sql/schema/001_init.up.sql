-- +goose Up
CREATE TABLE IF NOT EXISTS  classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_name VARCHAR(50) NOT NULL,
    year_of_study INT NOT NULL,
    graduation_year INT NOT NULL
);

CREATE TABLE IF NOT EXISTS  students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID REFERENCES classes (id) ON DELETE SET NULL,
    full_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS  teachers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS  subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(100) NOT NULL
);

-- Further, teacher_subject can be abbreviated as ts
CREATE TABLE IF NOT EXISTS  teacher_subject (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects (id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    UNIQUE (teacher_id, subject_id, class_id)
);

CREATE TABLE IF NOT EXISTS  status_code (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    abbreviation VARCHAR(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS  grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    ts_id UUID NOT NULL REFERENCES teacher_subject (id) ON DELETE CASCADE,
    status_code_id UUID REFERENCES status_code (id) ON DELETE CASCADE,
    score SMALLINT CHECK (score BETWEEN 2 AND 5),
    lesson_number SMALLINT NOT NULL CHECK (lesson_number >= 1),
    lesson_date DATE NOT NULL DEFAULT current_date,
    note TEXT,

    CONSTRAINT grade_or_status CHECK (
        (score IS NOT NULL AND status_code_id IS NULL)
        OR (score IS NULL AND status_code_id IS NOT NULL)
    ),

    CONSTRAINT unique_student_lesson
    UNIQUE (student_id, lesson_date, lesson_number)
);


CREATE TABLE IF NOT EXISTS  homework (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ts_id UUID NOT NULL REFERENCES teacher_subject (id) ON DELETE CASCADE,
    description_task TEXT NOT NULL,
    deadline_at DATE,
    assigned_at DATE NOT NULL DEFAULT current_date
);

INSERT INTO status_code (abbreviation) VALUES
('DEBT'),
('ABSENT'),
('SICK'),
('EXCUSE'),
('NA'),
('RETAKE'),
('EXEMPT');

-- +goose Down
DROP TABLE IF EXISTS homework;
DROP TABLE IF EXISTS grades;
DROP TABLE IF EXISTS status_code;
DROP TABLE IF EXISTS teacher_subject;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS teachers;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS classes;
