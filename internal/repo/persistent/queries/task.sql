-- name: CreateTask :exec
INSERT INTO tasks (id, user_id, title, description, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetTaskByID :one
SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1;

-- name: ListTasks :many
SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks
WHERE user_id = $1
  AND ($2::varchar IS NULL OR status = $2)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountTasks :one
SELECT COUNT(*) FROM tasks
WHERE user_id = $1
  AND ($2::varchar IS NULL OR status = $2);

-- name: UpdateTask :execrows
UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = $4
WHERE id = $5 AND user_id = $6;

-- name: DeleteTask :execrows
DELETE FROM tasks WHERE id = $1 AND user_id = $2;
