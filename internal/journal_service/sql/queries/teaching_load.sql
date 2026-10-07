-- name: CreateTeachingLoad :one
INSERT INTO teacher_subject (
    teacher_id,
    subject_id,
    class_id
) VALUES ($1, $2, $3)
RETURNING id, teacher_id, subject_id, class_id;

-- name: GetTeachingLoad :one
SELECT
    id,
    teacher_id,
    subject_id,
    class_id
FROM teacher_subject
WHERE id = $1;

-- name: GetTeachingID :one
SELECT id
FROM teacher_subject
WHERE 
    teacher_id = $1
    AND subject_id = $2
    AND class_id = $3;

-- name: ListTeachingLoad :many
-- name: ListTeachingLoad :many
SELECT
    id,
    teacher_id,
    subject_id,
    class_id
FROM teacher_subject
WHERE
    (sqlc.narg('teacher_id')::UUID IS NULL OR teacher_id = sqlc.narg('teacher_id'))
    AND (sqlc.narg('class_id')::UUID IS NULL OR class_id = sqlc.narg('class_id'));


-- name: UpdateTeachingLoad :one
UPDATE teacher_subject
SET
    teacher_id = coalesce(sqlc.narg('teacher_id'), teacher_id),
    subject_id = coalesce(sqlc.narg('subject_id'), subject_id),
    class_id = coalesce(sqlc.narg('class_id'), class_id)
WHERE id = $1
RETURNING id, teacher_id, subject_id, class_id;

-- name: DeleteTeachingLoad :exec
DELETE FROM teacher_subject
WHERE id = $1;
