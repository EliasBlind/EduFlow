-- name: CreateClass :one
INSERT INTO classes (
    class_name,
    year_of_study,
    graduation_year
) VALUES (
    $1, $2, $3
)
RETURNING id, class_name, year_of_study, graduation_year;

-- name: GetClass :one
SELECT
    id,
    class_name,
    year_of_study,
    graduation_year
FROM classes
WHERE id = $1
ORDER BY id
LIMIT 1;

-- name: ListClasses :many
SELECT DISTINCT
    id,
    class_name,
    year_of_study,
    graduation_year
FROM classes
ORDER BY
    year_of_study DESC,
    class_name ASC;

-- name: ListTeacherClasses :many
SELECT DISTINCT
    classes.id,
    classes.class_name,
    classes.year_of_study,
    classes.graduation_year
FROM classes
INNER JOIN teacher_subject ON classes.id = teacher_subject.class_id
WHERE teacher_subject.teacher_id = $1
ORDER BY
    classes.year_of_study DESC,
    classes.class_name ASC;

-- name: UpdateClass :one
UPDATE classes
SET
    class_name = coalesce(sqlc.narg('class_name'), class_name),
    year_of_study = coalesce(sqlc.narg('year_of_study'), year_of_study),
    graduation_year = coalesce(sqlc.narg('graduation_year'), graduation_year)
WHERE id = $1
RETURNING id, class_name, year_of_study, graduation_year;

-- name: DeleteClass :exec
DELETE FROM classes
WHERE id = $1;
