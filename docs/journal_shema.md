# Схема БД журнала

``` mermaid
erDiagram
    classes {
        UUID id PK "DEFAULT gen_random_uuid()"
        VARCHAR_50 class_name "NOT NULL"
        INT year_of_study "NOT NULL"
        INT graduation_year "NOT NULL"
    }

    students {
        UUID id PK "DEFAULT gen_random_uuid()"
        UUID class_id FK "REFERENCES classes (ON DELETE SET NULL)"
        VARCHAR_255 full_name "NOT NULL"
    }

    teachers {
        UUID id PK "DEFAULT gen_random_uuid()"
        VARCHAR_255 full_name "NOT NULL"
    }

    subjects {
        UUID id PK "DEFAULT gen_random_uuid()"
        VARCHAR_100 full_name "NOT NULL"
    }

    teacher_subject {
        UUID id PK "DEFAULT gen_random_uuid()"
        UUID teacher_id FK "REFERENCES teachers (ON DELETE CASCADE)"
        UUID subject_id FK "REFERENCES subjects (ON DELETE CASCADE)"
        UUID class_id FK "REFERENCES classes (ON DELETE CASCADE)"
    }

    status_code {
        UUID id PK "DEFAULT gen_random_uuid()"
        VARCHAR_6 abbreviation "NOT NULL"
    }

    grades {
        UUID id PK "DEFAULT gen_random_uuid()"
        UUID student_id FK "REFERENCES students (ON DELETE CASCADE)"
        UUID ts_id FK "REFERENCES teacher_subject (ON DELETE CASCADE)"
        UUID status_code_id FK "REFERENCES status_code (ON DELETE CASCADE)"
        SMALLINT score "CHECK (2..5)"
        SMALLINT lesson_number "NOT NULL CHECK (>=1)"
        DATE lesson_date "NOT NULL DEFAULT current_date"
        TEXT note
    }

    homework {
        UUID id PK "DEFAULT gen_random_uuid()"
        UUID ts_id FK "REFERENCES teacher_subject (ON DELETE CASCADE)"
        TEXT description_task "NOT NULL"
        DATE deadline_at
        DATE assigned_at "NOT NULL DEFAULT current_date"
    }

    %% Отношения между таблицами
    classes ||--o{ students : "class_id (SET NULL)"
    classes ||--o{ teacher_subject : "class_id (CASCADE)"
    teachers ||--o{ teacher_subject : "teacher_id (CASCADE)"
    subjects ||--o{ teacher_subject : "subject_id (CASCADE)"
    
    students ||--o{ grades : "student_id (CASCADE)"
    status_code ||--o{ grades : "status_code_id (CASCADE)"
    
    teacher_subject ||--o{ grades : "ts_id (CASCADE)"
    teacher_subject ||--o{ homework : "ts_id (CASCADE)"

```
