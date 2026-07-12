package v1

import (
	"context"
	"fmt"

	v1 "github.com/minhhoccode111/go-clean-template-gin/docs/proto/v1"
	grpcmw "github.com/minhhoccode111/go-clean-template-gin/internal/controller/grpc/middleware"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/grpc/v1/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *TranslationController) GetHistory(ctx context.Context, _ *v1.GetHistoryRequest) (*v1.GetHistoryResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	translationHistory, err := r.t.History(ctx, userID)
	if err != nil {
		r.l.Error(err, "grpc - v1 - GetHistory")

		return nil, fmt.Errorf("grpc - v1 - GetHistory: %w", err)
	}

	return response.NewTranslationHistory(translationHistory), nil
}
