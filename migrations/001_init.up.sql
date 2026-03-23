CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    class_name VARCHAR(50) NOT NULL,
    year_of_study INT NOT NULL,
    graduation_year INT NOT NULL
);

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    class_id UUID REFERENCES classes (id) ON DELETE SET NULL,
    full_name VARCHAR(255) NOT NULL
);

CREATE TABLE teachers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    full_name VARCHAR(255) NOT NULL
);

CREATE TABLE subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    full_name VARCHAR(100) NOT NULL
);

-- Further, teacher_subject can be abbreviated as ts
CREATE TABLE teacher_subject (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects (id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    UNIQUE (teacher_id, subject_id, class_id)
);

CREATE TABLE status_code (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    abbreviation VARCHAR(6) NOT NULL
);

CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
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


CREATE TABLE homework (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    description_task TEXT NOT NULL,
    deadline_at DATE,
    assigned_at DATE NOT NULL DEFAULT current_date
);

INSERT INTO grade_statuses (code) VALUES
('DEBT'),
('ABSENT'),
('SICK'),
('EXCUSE'),
('NA'),
('RETAKE'),
('EXEMPT')
