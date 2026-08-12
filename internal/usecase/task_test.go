package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/repo"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errRepoGeneric = errors.New("repository error")

const (
	testTaskID       = "task-id-123"
	testTaskUserID   = "user-id-123"
	testTaskTitle    = "My Task"
	testTaskOldTitle = "Old Title"
)

func newTaskUseCase(t *testing.T) (usecase.Task, *MockTaskRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mockRepo := NewMockTaskRepo(ctrl)
	useCase := task.New(mockRepo)

	return useCase, mockRepo
}

func TestTaskCreate(t *testing.T) {
	t.Parallel()

	t.Run("create success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)
		mockRepo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

		t2, err := uc.Create(context.Background(), testTaskUserID, testTaskTitle, "Task description")

		require.NoError(t, err)
		assert.NotEmpty(t, t2.ID)
		assert.Equal(t, testTaskTitle, t2.Title)
		assert.Equal(t, entity.TaskStatusTodo, t2.Status)
	})
}

func TestTaskGet(t *testing.T) {
	t.Parallel()

	expectedTask := entity.Task{
		ID:     testTaskID,
		UserID: testTaskUserID,
		Title:  testTaskTitle,
		Status: entity.TaskStatusTodo,
	}

	t.Run("get success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)
		mockRepo.EXPECT().GetByID(gomock.Any(), testTaskUserID, testTaskID).Return(expectedTask, nil)

		t2, err := uc.Get(context.Background(), testTaskUserID, testTaskID)

		require.NoError(t, err)
		assert.Equal(t, expectedTask, t2)
	})

	t.Run("get not found", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)
		mockRepo.EXPECT().GetByID(gomock.Any(), testTaskUserID, "missing-id").Return(entity.Task{}, entity.ErrTaskNotFound)

		_, err := uc.Get(context.Background(), testTaskUserID, "missing-id")

		require.ErrorIs(t, err, entity.ErrTaskNotFound)
	})
}

func TestTaskList(t *testing.T) {
	t.Parallel()

	task1 := entity.Task{ID: "task-1", UserID: testTaskUserID, Title: "Task 1", Status: entity.TaskStatusTodo}
	task2 := entity.Task{ID: "task-2", UserID: testTaskUserID, Title: "Task 2", Status: entity.TaskStatusInProgress}

	t.Run("list success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)
		mockRepo.EXPECT().List(gomock.Any(), testTaskUserID, gomock.Any()).Return([]entity.Task{task1, task2}, 2, nil)

		tasks, total, err := uc.List(context.Background(), testTaskUserID, nil, 10, 0)

		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, tasks, 2)
	})

	t.Run("list defaults", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)
		mockRepo.EXPECT().List(gomock.Any(), testTaskUserID, repo.TaskFilter{
			Status: nil,
			Limit:  uint64(10),
			Offset: uint64(0),
		}).Return([]entity.Task{task1, task2}, 2, nil)

		tasks, total, err := uc.List(context.Background(), testTaskUserID, nil, 0, -1)

		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, tasks, 2)
	})
}

func TestTaskUpdate(t *testing.T) {
	t.Parallel()

	t.Run("update success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)

		existingTask := entity.Task{
			ID:     testTaskID,
			UserID: testTaskUserID,
			Title:  testTaskOldTitle,
			Status: entity.TaskStatusTodo,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), testTaskUserID, testTaskID).Return(existingTask, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		updated, err := uc.Update(context.Background(), testTaskUserID, testTaskID, "New Title", "New description")

		require.NoError(t, err)
		assert.Equal(t, "New Title", updated.Title)
	})
}

func TestTaskTransition(t *testing.T) {
	t.Parallel()

	t.Run("transition valid", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)

		todoTask := entity.Task{
			ID:     testTaskID,
			UserID: testTaskUserID,
			Title:  testTaskTitle,
			Status: entity.TaskStatusTodo,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), testTaskUserID, testTaskID).Return(todoTask, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		updated, err := uc.Transition(context.Background(), testTaskUserID, testTaskID, entity.TaskStatusInProgress)

		require.NoError(t, err)
		assert.Equal(t, entity.TaskStatusInProgress, updated.Status)
	})

	t.Run("transition invalid", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)

		doneTask := entity.Task{
			ID:     "task-id-456",
			UserID: testTaskUserID,
			Title:  "Done Task",
			Status: entity.TaskStatusDone,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), testTaskUserID, "task-id-456").Return(doneTask, nil)

		_, err := uc.Transition(context.Background(), testTaskUserID, "task-id-456", entity.TaskStatusTodo)

		require.ErrorIs(t, err, entity.ErrInvalidTransition)
	})
}

func TestTaskDelete(t *testing.T) {
	t.Parallel()

	t.Run("delete success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)
		mockRepo.EXPECT().Delete(gomock.Any(), testTaskUserID, testTaskID).Return(nil)

		err := uc.Delete(context.Background(), testTaskUserID, testTaskID)

		require.NoError(t, err)
	})

	t.Run("delete not found", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo := newTaskUseCase(t)
		mockRepo.EXPECT().Delete(gomock.Any(), testTaskUserID, "missing-id").Return(entity.ErrTaskNotFound)

		err := uc.Delete(context.Background(), testTaskUserID, "missing-id")

		require.ErrorIs(t, err, entity.ErrTaskNotFound)
	})
}

func TestTaskCreate_RepoError(t *testing.T) {
	t.Parallel()

	uc, mockRepo := newTaskUseCase(t)

	mockRepo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(errRepoGeneric)

	_, err := uc.Create(context.Background(), testTaskUserID, "title", "desc")

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskGet_Forbidden(t *testing.T) {
	t.Parallel()

	uc, mockRepo := newTaskUseCase(t)

	mockRepo.EXPECT().GetByID(gomock.Any(), testTaskUserID, "task-id-999").Return(entity.Task{}, entity.ErrTaskForbidden)

	_, err := uc.Get(context.Background(), testTaskUserID, "task-id-999")

	require.Error(t, err)
	require.ErrorIs(t, err, entity.ErrTaskForbidden)
}

func TestTaskUpdate_RepoError(t *testing.T) {
	t.Parallel()

	uc, mockRepo := newTaskUseCase(t)

	existing := entity.Task{
		ID:     testTaskID,
		UserID: testTaskUserID,
		Title:  testTaskOldTitle,
		Status: entity.TaskStatusTodo,
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), testTaskUserID, testTaskID).Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errRepoGeneric)

	_, err := uc.Update(context.Background(), testTaskUserID, testTaskID, "New Title", "desc")

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskUpdate_NotFound(t *testing.T) {
	t.Parallel()

	uc, mockRepo := newTaskUseCase(t)

	mockRepo.EXPECT().GetByID(gomock.Any(), testTaskUserID, "missing-id").Return(entity.Task{}, entity.ErrTaskNotFound)

	_, err := uc.Update(context.Background(), testTaskUserID, "missing-id", "title", "desc")

	require.Error(t, err)
	require.ErrorIs(t, err, entity.ErrTaskNotFound)
}

func TestTaskTransition_UpdateError(t *testing.T) {
	t.Parallel()

	uc, mockRepo := newTaskUseCase(t)

	todoTask := entity.Task{
		ID:     testTaskID,
		UserID: testTaskUserID,
		Title:  testTaskTitle,
		Status: entity.TaskStatusTodo,
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), testTaskUserID, testTaskID).Return(todoTask, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errRepoGeneric)

	_, err := uc.Transition(context.Background(), testTaskUserID, testTaskID, entity.TaskStatusInProgress)

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskDelete_GenericError(t *testing.T) {
	t.Parallel()

	uc, mockRepo := newTaskUseCase(t)

	mockRepo.EXPECT().Delete(gomock.Any(), testTaskUserID, testTaskID).Return(errRepoGeneric)

	err := uc.Delete(context.Background(), testTaskUserID, testTaskID)

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskList_RepoError(t *testing.T) {
	t.Parallel()

	uc, mockRepo := newTaskUseCase(t)

	mockRepo.EXPECT().
		List(gomock.Any(), testTaskUserID, repo.TaskFilter{Limit: uint64(10), Offset: uint64(0)}).
		Return(nil, 0, errRepoGeneric)

	_, _, err := uc.List(context.Background(), testTaskUserID, nil, 10, 0)

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}

func TestTaskTransition_NotFound(t *testing.T) {
	t.Parallel()

	uc, mockRepo := newTaskUseCase(t)

	mockRepo.EXPECT().
		GetByID(gomock.Any(), testTaskUserID, testTaskID).
		Return(entity.Task{}, entity.ErrTaskNotFound)

	_, err := uc.Transition(context.Background(), testTaskUserID, testTaskID, entity.TaskStatusInProgress)

	require.Error(t, err)
	require.ErrorIs(t, err, entity.ErrTaskNotFound)
}
