package v1

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/nats/nats_rpc/server"
	"github.com/nats-io/nats.go"
)

type natsTaskCreateData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type natsTaskGetData struct {
	ID string `json:"id"`
}

type natsTaskListData struct {
	Status string `json:"status"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type natsTaskUpdateData struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type natsTaskTransitionData struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type natsTaskDeleteData struct {
	ID string `json:"id"`
}

func (r *V1) createTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		var req struct {
			Token string             `json:"token"`
			Data  natsTaskCreateData `json:"data"`
		}

		if err := json.Unmarshal(msg.Data, &req); err != nil {
			r.l.Error(err, "nats_rpc - V1 - createTask")

			return nil, fmt.Errorf("nats_rpc - V1 - createTask - json.Unmarshal: %w", err)
		}

		userID, err := r.j.ParseToken(req.Token)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - createTask - invalid token: %w", err)
		}

		task, err := r.tk.Create(context.Background(), userID, req.Data.Title, req.Data.Description)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - createTask")

			return nil, fmt.Errorf("nats_rpc - V1 - createTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) getTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		var req struct {
			Token string          `json:"token"`
			Data  natsTaskGetData `json:"data"`
		}

		if err := json.Unmarshal(msg.Data, &req); err != nil {
			r.l.Error(err, "nats_rpc - V1 - getTask")

			return nil, fmt.Errorf("nats_rpc - V1 - getTask - json.Unmarshal: %w", err)
		}

		userID, err := r.j.ParseToken(req.Token)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - getTask - invalid token: %w", err)
		}

		task, err := r.tk.Get(context.Background(), userID, req.Data.ID)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - getTask")

			return nil, fmt.Errorf("nats_rpc - V1 - getTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) listTasks() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		var req struct {
			Token string           `json:"token"`
			Data  natsTaskListData `json:"data"`
		}

		if err := json.Unmarshal(msg.Data, &req); err != nil {
			r.l.Error(err, "nats_rpc - V1 - listTasks")

			return nil, fmt.Errorf("nats_rpc - V1 - listTasks - json.Unmarshal: %w", err)
		}

		userID, err := r.j.ParseToken(req.Token)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - listTasks - invalid token: %w", err)
		}

		var status *entity.TaskStatus
		if req.Data.Status != "" {
			s := entity.TaskStatus(req.Data.Status)
			status = &s
		}

		tasks, total, err := r.tk.List(context.Background(), userID, status, req.Data.Limit, req.Data.Offset)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - listTasks")

			return nil, fmt.Errorf("nats_rpc - V1 - listTasks: %w", err)
		}

		return map[string]any{"tasks": tasks, "total": total}, nil
	}
}

func (r *V1) updateTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		var req struct {
			Token string            `json:"token"`
			Data  natsTaskUpdateData `json:"data"`
		}

		if err := json.Unmarshal(msg.Data, &req); err != nil {
			r.l.Error(err, "nats_rpc - V1 - updateTask")

			return nil, fmt.Errorf("nats_rpc - V1 - updateTask - json.Unmarshal: %w", err)
		}

		userID, err := r.j.ParseToken(req.Token)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - updateTask - invalid token: %w", err)
		}

		task, err := r.tk.Update(context.Background(), userID, req.Data.ID, req.Data.Title, req.Data.Description)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - updateTask")

			return nil, fmt.Errorf("nats_rpc - V1 - updateTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) transitionTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		var req struct {
			Token string                `json:"token"`
			Data  natsTaskTransitionData `json:"data"`
		}

		if err := json.Unmarshal(msg.Data, &req); err != nil {
			r.l.Error(err, "nats_rpc - V1 - transitionTask")

			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask - json.Unmarshal: %w", err)
		}

		userID, err := r.j.ParseToken(req.Token)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask - invalid token: %w", err)
		}

		task, err := r.tk.Transition(context.Background(), userID, req.Data.ID, entity.TaskStatus(req.Data.Status))
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - transitionTask")

			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) deleteTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		var req struct {
			Token string            `json:"token"`
			Data  natsTaskDeleteData `json:"data"`
		}

		if err := json.Unmarshal(msg.Data, &req); err != nil {
			r.l.Error(err, "nats_rpc - V1 - deleteTask")

			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask - json.Unmarshal: %w", err)
		}

		userID, err := r.j.ParseToken(req.Token)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask - invalid token: %w", err)
		}

		err = r.tk.Delete(context.Background(), userID, req.Data.ID)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - deleteTask")

			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask: %w", err)
		}

		return nil, nil
	}
}
