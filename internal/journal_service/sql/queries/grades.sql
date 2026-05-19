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
    ts.subject_id = $1
    AND ts.class_id = $2
    AND (sqlc.narg('student_id')::uuid IS NULL OR g.student_id = sqlc.narg('student_id')::uuid)
ORDER BY g.lesson_date;

-- name: UpdateGrades :one
UPDATE grades
SET
    status_code_id = $2,
    score = $3,
    note = coalesce(sqlc.narg('note'), note)
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
