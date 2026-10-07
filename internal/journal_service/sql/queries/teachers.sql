-- name: CreateTeacher :one
INSERT INTO teachers (
    id,
    full_name
) VALUES (
    COALESCE(sqlc.narg('id')::uuid, gen_random_uuid()), 
    sqlc.arg('full_name')
)
RETURNING id, full_name;

-- name: ListTeachers :many
SELECT
    id,
    full_name
FROM teachers
ORDER BY full_name
LIMIT
    $1
    OFFSET
    $2;

-- name: UpdateTeacher :one
UPDATE teachers
SET
    full_name = $2
WHERE id = $1
RETURNING id, full_name;

-- name: DeleteTeacher :exec
DELETE FROM teachers
WHERE id = $1;
