package entity

import "testing"

func TestTask_Transition(t *testing.T) {
	tests := []struct {
		name     string
		initial  TaskStatus
		new      TaskStatus
		wantErr  bool
		wantS    TaskStatus
	}{
		{"todo to in_progress", TaskStatusTodo, TaskStatusInProgress, false, TaskStatusInProgress},
		{"todo to done", TaskStatusTodo, TaskStatusDone, true, TaskStatusTodo},
		{"in_progress to done", TaskStatusInProgress, TaskStatusDone, false, TaskStatusDone},
		{"in_progress to todo", TaskStatusInProgress, TaskStatusTodo, false, TaskStatusTodo},
		{"done to todo", TaskStatusDone, TaskStatusTodo, true, TaskStatusDone},
		{"done to in_progress", TaskStatusDone, TaskStatusInProgress, true, TaskStatusDone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{Status: tt.initial}

			err := task.Transition(tt.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transition() error = %v, wantErr %v", err, tt.wantErr)
			}

			if task.Status != tt.wantS {
				t.Errorf("Transition() status = %v, want %v", task.Status, tt.wantS)
			}
		})
	}
}
