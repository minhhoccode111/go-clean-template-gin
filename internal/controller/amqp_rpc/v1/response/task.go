package response

import "github.com/minhhoccode111/go-clean-template-gin/internal/entity"

// TaskList -.
type TaskList struct {
	Tasks []entity.Task `json:"tasks"`
	Total int           `json:"total"`
}

// DeleteStatus -.
type DeleteStatus struct {
	Status string `json:"status"`
}
