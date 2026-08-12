package v1

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/nats_rpc/v1/request"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/nats_rpc/v1/response"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/nats/nats_rpc/server"
	"github.com/nats-io/nats.go"
)

func (r *V1) createTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, rawData, err := extractUserID(msg, r.j)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - createTask")

			return nil, fmt.Errorf("nats_rpc - V1 - createTask - extractUserID: %w", err)
		}

		var reqData request.CreateTask

		if err := json.Unmarshal(rawData, &reqData); err != nil {
			r.l.Error(err, "nats_rpc - V1 - createTask")

			return nil, fmt.Errorf("nats_rpc - V1 - createTask - json.Unmarshal: %w", err)
		}

		task, err := r.tk.Create(context.Background(), userID, reqData.Title, reqData.Description)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - createTask")

			return nil, fmt.Errorf("nats_rpc - V1 - createTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) getTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, rawData, err := extractUserID(msg, r.j)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - getTask")

			return nil, fmt.Errorf("nats_rpc - V1 - getTask - extractUserID: %w", err)
		}

		var reqData request.GetTask

		if err := json.Unmarshal(rawData, &reqData); err != nil {
			r.l.Error(err, "nats_rpc - V1 - getTask")

			return nil, fmt.Errorf("nats_rpc - V1 - getTask - json.Unmarshal: %w", err)
		}

		task, err := r.tk.Get(context.Background(), userID, reqData.ID)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - getTask")

			return nil, fmt.Errorf("nats_rpc - V1 - getTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) listTasks() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, rawData, err := extractUserID(msg, r.j)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - listTasks")

			return nil, fmt.Errorf("nats_rpc - V1 - listTasks - extractUserID: %w", err)
		}

		var reqData request.ListTasks

		if err := json.Unmarshal(rawData, &reqData); err != nil {
			r.l.Error(err, "nats_rpc - V1 - listTasks")

			return nil, fmt.Errorf("nats_rpc - V1 - listTasks - json.Unmarshal: %w", err)
		}

		var status *entity.TaskStatus

		if reqData.Status != "" {
			s := entity.TaskStatus(reqData.Status)
			status = &s
		}

		tasks, total, err := r.tk.List(context.Background(), userID, status, reqData.Limit, reqData.Offset)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - listTasks")

			return nil, fmt.Errorf("nats_rpc - V1 - listTasks: %w", err)
		}

		return response.TaskList{Tasks: tasks, Total: total}, nil
	}
}

func (r *V1) updateTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, rawData, err := extractUserID(msg, r.j)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - updateTask")

			return nil, fmt.Errorf("nats_rpc - V1 - updateTask - extractUserID: %w", err)
		}

		var reqData request.UpdateTask

		if err := json.Unmarshal(rawData, &reqData); err != nil {
			r.l.Error(err, "nats_rpc - V1 - updateTask")

			return nil, fmt.Errorf("nats_rpc - V1 - updateTask - json.Unmarshal: %w", err)
		}

		task, err := r.tk.Update(context.Background(), userID, reqData.ID, reqData.Title, reqData.Description)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - updateTask")

			return nil, fmt.Errorf("nats_rpc - V1 - updateTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) transitionTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, rawData, err := extractUserID(msg, r.j)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - transitionTask")

			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask - extractUserID: %w", err)
		}

		var reqData request.TransitionTask

		if err := json.Unmarshal(rawData, &reqData); err != nil {
			r.l.Error(err, "nats_rpc - V1 - transitionTask")

			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask - json.Unmarshal: %w", err)
		}

		task, err := r.tk.Transition(context.Background(), userID, reqData.ID, entity.TaskStatus(reqData.Status))
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - transitionTask")

			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) deleteTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, rawData, err := extractUserID(msg, r.j)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - deleteTask")

			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask - extractUserID: %w", err)
		}

		var reqData request.DeleteTask

		if err := json.Unmarshal(rawData, &reqData); err != nil {
			r.l.Error(err, "nats_rpc - V1 - deleteTask")

			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask - json.Unmarshal: %w", err)
		}

		err = r.tk.Delete(context.Background(), userID, reqData.ID)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - deleteTask")

			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask: %w", err)
		}

		return response.DeleteStatus{Status: "deleted"}, nil
	}
}
