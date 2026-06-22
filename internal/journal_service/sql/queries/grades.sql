-- name: RecordGrade :one
INSERT INTO grades (
    student_id,
    ts_id,
    status_code_id,
    score,
    lesson_number,
    lesson_date,
    note
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT ON CONSTRAINT unique_student_lesson 
DO UPDATE SET
    status_code_id = EXCLUDED.status_code_id,
    score = EXCLUDED.score,
    note = EXCLUDED.note,
    lesson_date = EXCLUDED.lesson_date
RETURNING
    id,
    student_id,
    ts_id,
    status_code_id,
    score,
    lesson_number,
    lesson_date,
    note;

-- name: GetGrades :many
SELECT
    g.id,
    ts.subject_id,
    g.student_id,
    ts.class_id,
    g.status_code_id,
    g.score,
    g.lesson_number,
    g.lesson_date,
    g.note
FROM grades AS g
INNER JOIN teacher_subject AS ts ON g.ts_id = ts.id
INNER JOIN students AS st ON g.student_id = st.id
WHERE
    (sqlc.narg('subject_id')::uuid IS NULL OR ts.subject_id = sqlc.narg('subject_id')::uuid)
    AND (sqlc.narg('class_id')::uuid IS NULL OR ts.class_id = sqlc.narg('class_id')::uuid)
    AND (sqlc.narg('student_id')::uuid IS NULL OR g.student_id = sqlc.narg('student_id')::uuid)
ORDER BY g.lesson_date;

-- name: UpdateGrades :one
UPDATE grades
SET
    status_code_id = $2,
    score = $3,
    note = $4 -- Теперь NULL полностью затрет старую заметку
WHERE id = $1
RETURNING
    id,
    student_id,
    ts_id,
    status_code_id,
    score,
    lesson_number,
    lesson_date,
    note;

-- name: DeleteGrade :exec
DELETE FROM grades
WHERE id = $1;
