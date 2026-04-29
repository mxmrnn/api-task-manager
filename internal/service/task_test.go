package service_test

import (
	"async-api-task-manager/internal/model"
	"async-api-task-manager/internal/service"
	"async-api-task-manager/internal/storage/postgres/repository"
	"async-api-task-manager/internal/transport/http/dto"
	"async-api-task-manager/internal/worker"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) List(ctx context.Context, filter dto.Filter) ([]model.Task, error) {
	args := m.Called(ctx, filter)

	var tasks []model.Task
	if arg := args.Get(0); arg != nil {
		tasks = arg.([]model.Task)
	}

	return tasks, args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

var _ service.TaskRepository = (*MockTaskRepository)(nil)

func TestTaskService_GetTaskByID(t *testing.T) {
	tests := []struct {
		name      string
		taskID    uuid.UUID
		mockSetup func(*MockTaskRepository)
		wantErr   error
		wantTitle string
	}{
		{
			name:   "success get by id task",
			taskID: uuid.New(),
			mockSetup: func(m *MockTaskRepository) {
				task := &model.Task{
					ID:          uuid.New(),
					Title:       "Test task",
					Description: "Test description",
					Status:      model.TaskStatusTodo,
				}
				m.On("GetByID", mock.Anything, mock.Anything).Return(task, nil)
			},
			wantErr:   nil,
			wantTitle: "Test task",
		},
		{
			name:   "task not found",
			taskID: uuid.New(),
			mockSetup: func(m *MockTaskRepository) {
				m.On("GetByID", mock.Anything, mock.Anything).Return(nil, repository.ErrTaskNotFound)
			},
			wantErr: service.ErrTaskNotFound,
		},
		{
			name:   "else err of repository",
			taskID: uuid.New(),
			mockSetup: func(m *MockTaskRepository) {
				m.On("GetByID", mock.Anything, mock.Anything).Return(nil, errors.New("db connection error"))
			},
			wantErr: errors.New("db connection error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTaskRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewTaskService(mockRepo, &worker.Worker{})

			// Act
			resp, err := svc.GetTaskByID(context.Background(), tt.taskID)

			// Assert
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr) || err.Error() == tt.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.ID)
			assert.Equal(t, tt.wantTitle, resp.Title)
		})
	}
}
