-- name: CreateStatusCode :one
INSERT INTO status_code (
    abbreviation
) VALUES (
    $1
)
RETURNING id, abbreviation;

-- name: UpdateStatusCode :one
UPDATE status_code
SET
    abbreviation = $2
WHERE id = $1
RETURNING id, abbreviation;

-- name: ListStatusCode :many
SELECT
    id,
    abbreviation
FROM status_code
ORDER BY abbreviation ASC;

-- name: DeleteStatusCode :exec
DELETE FROM status_code
WHERE id = $1;
