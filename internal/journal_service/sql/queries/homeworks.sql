-- name: RecordHomework :one
INSERT INTO homework (
    ts_id,
    description_task,
    deadline_at,
    assigned_at
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, ts_id, description_task, deadline_at, assigned_at;

-- name: UpdateHomework :one
UPDATE homework
SET
    description_task = coalesce(sqlc.narg('description_task'), description_task),
    deadline_at = coalesce(sqlc.narg('deadline_at'), deadline_at),
    assigned_at = coalesce(sqlc.narg('assigned_at'), assigned_at)
WHERE
    id = $1
RETURNING id, ts_id, description_task, deadline_at, assigned_at;

-- name: ListHomeworks :many
SELECT
    h.id,
    ts.teacher_id,
    ts.subject_id,
    ts.class_id,
    h.description_task,
    h.deadline_at,
    h.assigned_at
FROM homework AS h
INNER JOIN teacher_subject AS ts ON homework.ts_id = teacher_subject.id
WHERE
    ts.class_id = $1
    AND ts.subject_id = $2
ORDER BY h.assigned_at DESC;

-- name: DeleteHomeworks :exec
DELETE FROM homework
WHERE id = $1;
