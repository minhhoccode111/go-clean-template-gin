package translation

import (
	"context"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/translation"

type tracedRepo struct {
	next repo.TranslationRepo
}

func newTraced(next repo.TranslationRepo) repo.TranslationRepo {
	return &tracedRepo{next: next}
}

func (r *tracedRepo) Store(ctx context.Context, userID string, t entity.Translation) error {
	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "TranslationRepo.Store",
		trace.WithAttributes(
			attribute.String("user.id", userID),
			attribute.String("translation.source", t.Source),
			attribute.String("translation.destination", t.Destination),
		),
	)

	err := r.next.Store(ctx, userID, t)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return err
}

func (r *tracedRepo) GetHistory(ctx context.Context, userID string) ([]entity.Translation, error) {
	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "TranslationRepo.GetHistory",
		trace.WithAttributes(attribute.String("user.id", userID)),
	)

	result, err := r.next.GetHistory(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return result, err
}
