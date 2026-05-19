-- name: CreateSubject :one
INSERT INTO subjects (
    full_name
) VALUES ($1)
RETURNING id, full_name;

-- name: UpdateSubject :one
UPDATE subjects
SET
    full_name = $2
WHERE id = $1
RETURNING id, full_name;

-- name: DeleteSubject :exec
DELETE FROM subjects
WHERE id = $1;

-- name: ListSubjects :many
SELECT
    id,
    full_name
FROM subjects
ORDER BY full_name;
