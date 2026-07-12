package v1

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/rabbitmq/rmq_rpc/server"
	amqp "github.com/rabbitmq/amqp091-go"
)

type translateReqData struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Original    string `json:"original"`
}

func (r *V1) getHistory() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		var req struct {
			Token string `json:"token"`
		}

		if err := json.Unmarshal(d.Body, &req); err != nil {
			r.l.Error(err, "amqp_rpc - V1 - getHistory")

			return nil, fmt.Errorf("amqp_rpc - V1 - getHistory - json.Unmarshal: %w", err)
		}

		userID, err := r.j.ParseToken(req.Token)
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - V1 - getHistory - invalid token: %w", err)
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
		var req struct {
			Token string           `json:"token"`
			Data  translateReqData `json:"data"`
		}

		if err := json.Unmarshal(d.Body, &req); err != nil {
			r.l.Error(err, "amqp_rpc - V1 - translate")

			return nil, fmt.Errorf("amqp_rpc - V1 - translate - json.Unmarshal: %w", err)
		}

		userID, err := r.j.ParseToken(req.Token)
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - V1 - translate - invalid token: %w", err)
		}

		translation, err := r.t.Translate(context.Background(), userID, entity.Translation{
			Source:      req.Data.Source,
			Destination: req.Data.Destination,
			Original:    req.Data.Original,
		})
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - translate")

			return nil, fmt.Errorf("amqp_rpc - V1 - translate: %w", err)
		}

		return translation, nil
	}
}
