package v1

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/amqp_rpc/v1/request"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/rabbitmq/rmq_rpc/server"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *V1) getHistory() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, _, err := extractUserID(d, r.j)
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - getHistory")

			return nil, fmt.Errorf("amqp_rpc - V1 - getHistory - extractUserID: %w", err)
		}

		translationHistory, err := r.t.History(context.Background(), userID)
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - getHistory")

			return nil, fmt.Errorf("amqp_rpc - V1 - getHistory: %w", err)
		}

		return translationHistory, nil
	}
}

func (r *V1) translate() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, rawData, err := extractUserID(d, r.j)
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - translate")

			return nil, fmt.Errorf("amqp_rpc - V1 - translate - extractUserID: %w", err)
		}

		var reqData request.Translate

		if err := json.Unmarshal(rawData, &reqData); err != nil {
			r.l.Error(err, "amqp_rpc - V1 - translate")

			return nil, fmt.Errorf("amqp_rpc - V1 - translate - json.Unmarshal: %w", err)
		}

		translation, err := r.t.Translate(context.Background(), userID, entity.Translation{
			Source:      reqData.Source,
			Destination: reqData.Destination,
			Original:    reqData.Original,
		})
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - translate")

			return nil, fmt.Errorf("amqp_rpc - V1 - translate: %w", err)
		}

		return translation, nil
	}
}
