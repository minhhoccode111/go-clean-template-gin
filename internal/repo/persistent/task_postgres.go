package persistent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/sqlc"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/postgres"
)

// toInt32 saturates a uint64 to int32 instead of overflowing (gosec G115).
func toInt32(v uint64) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}

	return int32(v)
}

// toInt64 saturates a uint64 to int64 instead of overflowing (gosec G115).
func toInt64(v uint64) int64 {
	if v > math.MaxInt64 {
		return math.MaxInt64
	}

	return int64(v)
}

// TaskRepo implements repo.TaskRepo using sqlc.
type TaskRepo struct {
	*postgres.Postgres
	queries *sqlc.Queries
}

// NewTaskRepo returns a Task repository instrumented with OpenTelemetry tracing spans.
func NewTaskRepo(pg *postgres.Postgres) repo.TaskRepo {
	return newTracedTask(&TaskRepo{
		Postgres: pg,
		queries:  sqlc.New(pg.Pool),
	})
}

// Store -.
func (r *TaskRepo) Store(ctx context.Context, task *entity.Task) error {
	err := r.queries.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		CreatedAt:   pgTimestamptz(task.CreatedAt),
		UpdatedAt:   pgTimestamptz(task.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("TaskRepo - Store - r.queries.CreateTask: %w", err)
	}

	return nil
}

// GetByID -.
func (r *TaskRepo) GetByID(ctx context.Context, userID, taskID string) (entity.Task, error) {
	row, err := r.queries.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Task{}, entity.ErrTaskNotFound
		}

		return entity.Task{}, fmt.Errorf("TaskRepo - GetByID - r.queries.GetTaskByID: %w", err)
	}

	if row.UserID != userID {
		return entity.Task{}, entity.ErrTaskForbidden
	}

	return entity.Task{
		ID:          row.ID,
		UserID:      row.UserID,
		Title:       row.Title,
		Description: row.Description,
		Status:      entity.TaskStatus(row.Status),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}, nil
}

// List -.
func (r *TaskRepo) List(ctx context.Context, userID string, filter repo.TaskFilter) ([]entity.Task, int, error) {
	var statusStr string
	if filter.Status != nil {
		statusStr = string(*filter.Status)
	}

	total, err := r.queries.CountTasks(ctx, sqlc.CountTasksParams{
		UserID:  userID,
		Column2: statusStr,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - r.queries.CountTasks: %w", err)
	}

	rows, err := r.queries.ListTasks(ctx, sqlc.ListTasksParams{
		UserID:  userID,
		Column2: statusStr,
		Limit:   toInt32(filter.Limit),
		Offset:  toInt32(filter.Offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - r.queries.ListTasks: %w", err)
	}

	tasks := make([]entity.Task, 0, len(rows))
	for i := range rows {
		tasks = append(tasks, entity.Task{
			ID:          rows[i].ID,
			UserID:      rows[i].UserID,
			Title:       rows[i].Title,
			Description: rows[i].Description,
			Status:      entity.TaskStatus(rows[i].Status),
			CreatedAt:   rows[i].CreatedAt.Time,
			UpdatedAt:   rows[i].UpdatedAt.Time,
		})
	}

	return tasks, int(total), nil
}

// Update -.
func (r *TaskRepo) Update(ctx context.Context, task *entity.Task) error {
	rowsAffected, err := r.queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		UpdatedAt:   pgTimestamptz(task.UpdatedAt),
		ID:          task.ID,
		UserID:      task.UserID,
	})
	if err != nil {
		return fmt.Errorf("TaskRepo - Update - r.queries.UpdateTask: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrTaskNotFound
	}

	return nil
}

// Delete -.
func (r *TaskRepo) Delete(ctx context.Context, userID, taskID string) error {
	rowsAffected, err := r.queries.DeleteTask(ctx, sqlc.DeleteTaskParams{
		ID:     taskID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("TaskRepo - Delete - r.queries.DeleteTask: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrTaskNotFound
	}

	return nil
}
