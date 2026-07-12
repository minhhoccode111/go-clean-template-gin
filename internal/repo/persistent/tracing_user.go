package persistent

import (
	"context"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerNameUser = "github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/user"

type tracedUserRepo struct {
	next repo.UserRepo
}

func newTracedUser(next repo.UserRepo) repo.UserRepo {
	return &tracedUserRepo{next: next}
}

func (r *tracedUserRepo) Store(ctx context.Context, user *entity.User) error {
	ctx, span := otel.Tracer(_tracerNameUser).Start(
		ctx, "UserRepo.Store",
		trace.WithAttributes(
			attribute.String("user.id", user.ID),
			attribute.String("user.email", user.Email),
		),
	)

	err := r.next.Store(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return err
}

func (r *tracedUserRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	ctx, span := otel.Tracer(_tracerNameUser).Start(
		ctx, "UserRepo.GetByID",
		trace.WithAttributes(attribute.String("user.id", id)),
	)

	result, err := r.next.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return result, err
}

func (r *tracedUserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	ctx, span := otel.Tracer(_tracerNameUser).Start(
		ctx, "UserRepo.GetByEmail",
		trace.WithAttributes(attribute.String("user.email", email)),
	)

	result, err := r.next.GetByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return result, err
}
