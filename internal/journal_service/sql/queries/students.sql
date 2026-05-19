-- name: CreateStudent :one
INSERT INTO students (
    id,
    class_id,
    full_name
) VALUES (
    COALESCE(sqlc.narg('id')::uuid, gen_random_uuid()), 
    sqlc.arg('class_id'), 
    sqlc.arg('full_name')
)
RETURNING id, class_id, full_name;

-- name: GetStudent :one
SELECT
    id,
    class_id,
    full_name
FROM students
WHERE id = sqlc.arg('id')
ORDER BY full_name
LIMIT 1;

-- name: ListStudents :many
SELECT
    id,
    class_id,
    full_name
FROM students
WHERE class_id = sqlc.arg('class_id')
ORDER BY full_name
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: UpdateStudent :one
UPDATE students
SET
    class_id = COALESCE(sqlc.narg('class_id'), class_id),
    full_name = COALESCE(sqlc.narg('full_name'), full_name)
WHERE
    id = sqlc.arg('id')
RETURNING id, class_id, full_name;

-- name: DeleteStudent :exec
DELETE FROM students
WHERE id = sqlc.arg('id');
