package v1

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/nats_rpc/v1/request"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/nats/nats_rpc/server"
	"github.com/nats-io/nats.go"
)

func (r *V1) getHistory() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, _, err := extractUserID(msg, r.j)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - getHistory")

			return nil, fmt.Errorf("nats_rpc - V1 - getHistory - extractUserID: %w", err)
		}

		translationHistory, err := r.t.History(context.Background(), userID)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - getHistory")

			return nil, fmt.Errorf("nats_rpc - V1 - getHistory: %w", err)
		}

		return translationHistory, nil
	}
}

func (r *V1) translate() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, rawData, err := extractUserID(msg, r.j)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - translate")

			return nil, fmt.Errorf("nats_rpc - V1 - translate - extractUserID: %w", err)
		}

		var reqData request.Translate

		if err := json.Unmarshal(rawData, &reqData); err != nil {
			r.l.Error(err, "nats_rpc - V1 - translate")

			return nil, fmt.Errorf("nats_rpc - V1 - translate - json.Unmarshal: %w", err)
		}

		translation, err := r.t.Translate(context.Background(), userID, entity.Translation{
			Source:      reqData.Source,
			Destination: reqData.Destination,
			Original:    reqData.Original,
		})
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - translate")

			return nil, fmt.Errorf("nats_rpc - V1 - translate: %w", err)
		}

		return translation, nil
	}
}
