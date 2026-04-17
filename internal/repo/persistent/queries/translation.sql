-- name: GetHistory :many
SELECT source, destination, original, translation FROM histories;

-- name: Store :exec
INSERT INTO histories (source, destination, original, translation)
VALUES ($1, $2, $3, $4);
