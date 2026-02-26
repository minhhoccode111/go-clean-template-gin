-- name: GetHistory :many
SELECT source, destination, original, translation FROM history;

-- name: Store :exec
INSERT INTO history (source, destination, original, translation)
VALUES ($1, $2, $3, $4);
