package task

import (
	"context"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/task"

type tracedRepo struct {
	next repo.TaskRepo
}

func newTraced(next repo.TaskRepo) repo.TaskRepo {
	return &tracedRepo{next: next}
}

func (r *tracedRepo) Store(ctx context.Context, task *entity.Task) error {
	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "TaskRepo.Store",
		trace.WithAttributes(
			attribute.String("task.id", task.ID),
			attribute.String("user.id", task.UserID),
		),
	)

	err := r.next.Store(ctx, task)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return err
}

func (r *tracedRepo) GetByID(ctx context.Context, userID, taskID string) (entity.Task, error) {
	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "TaskRepo.GetByID",
		trace.WithAttributes(
			attribute.String("user.id", userID),
			attribute.String("task.id", taskID),
		),
	)

	result, err := r.next.GetByID(ctx, userID, taskID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return result, err
}

func (r *tracedRepo) List(ctx context.Context, userID string, filter repo.TaskFilter) ([]entity.Task, int, error) {
	attrs := []attribute.KeyValue{
		attribute.String("user.id", userID),
		attribute.Int64("task.limit", toInt64(filter.Limit)),
		attribute.Int64("task.offset", toInt64(filter.Offset)),
	}
	if filter.Status != nil {
		attrs = append(attrs, attribute.String("task.status", string(*filter.Status)))
	}

	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "TaskRepo.List",
		trace.WithAttributes(attrs...),
	)

	tasks, total, err := r.next.List(ctx, userID, filter)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return tasks, total, err
}

func (r *tracedRepo) Update(ctx context.Context, task *entity.Task) error {
	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "TaskRepo.Update",
		trace.WithAttributes(
			attribute.String("task.id", task.ID),
			attribute.String("user.id", task.UserID),
		),
	)

	err := r.next.Update(ctx, task)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return err
}

func (r *tracedRepo) Delete(ctx context.Context, userID, taskID string) error {
	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "TaskRepo.Delete",
		trace.WithAttributes(
			attribute.String("user.id", userID),
			attribute.String("task.id", taskID),
		),
	)

	err := r.next.Delete(ctx, userID, taskID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return err
}
