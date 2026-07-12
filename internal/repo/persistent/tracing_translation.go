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

const _tracerNameTranslation = "github.com/minhhoccode111/go-clean-template-gin/internal/repo/persistent/translation"

type tracedTranslationRepo struct {
	next repo.TranslationRepo
}

func newTracedTranslation(next repo.TranslationRepo) repo.TranslationRepo {
	return &tracedTranslationRepo{next: next}
}

func (r *tracedTranslationRepo) Store(ctx context.Context, userID string, t entity.Translation) error {
	ctx, span := otel.Tracer(_tracerNameTranslation).Start(
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

func (r *tracedTranslationRepo) GetHistory(ctx context.Context, userID string) ([]entity.Translation, error) {
	ctx, span := otel.Tracer(_tracerNameTranslation).Start(
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
