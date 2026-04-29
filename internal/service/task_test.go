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
	"time"

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

func TestTaskService_ListTasks(t *testing.T) {
	tests := []struct {
		name      string
		filter    dto.Filter
		mockSetup func(*MockTaskRepository)
		wantErr   error
		wantCount int
	}{
		{
			name:   "success list tasks without filter",
			filter: dto.Filter{},
			mockSetup: func(m *MockTaskRepository) {
				tasks := []model.Task{
					{ID: uuid.New(), Title: "Task 1", Status: model.TaskStatusTodo, Description: "Description 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
					{ID: uuid.New(), Title: "Task 2", Status: model.TaskStatusDone, Description: "Description 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.Filter")).Return(tasks, nil)
			},
			wantErr:   nil,
			wantCount: 2,
		},
		{
			name:   "success list tasks with status filter",
			filter: dto.Filter{Status: ptr(t, model.TaskStatusInProgress)},
			mockSetup: func(m *MockTaskRepository) {
				tasks := []model.Task{{ID: uuid.New(), Title: "In Progress Task", Status: model.TaskStatusTodo}}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.Filter")).Return(tasks, nil)
			},
			wantErr:   nil,
			wantCount: 1,
		},
		{
			name: "err filter with not valid status",
			filter: dto.Filter{Status: func() *model.TaskStatus {
				s := model.TaskStatus("invalid_status")
				return &s
			}()},
			mockSetup: func(m *MockTaskRepository) {
			},
			wantErr: service.ErrInvalidTaskStatus,
		},
		{
			name:   "err from repository",
			filter: dto.Filter{},
			mockSetup: func(m *MockTaskRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.Filter")).
					Return(nil, errors.New("database connection error"))
			},
			wantErr: errors.New("database connection error"),
		},
		{
			name:   "empty list of tasks",
			filter: dto.Filter{},
			mockSetup: func(m *MockTaskRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.Filter")).
					Return([]model.Task{}, nil)
			},
			wantErr:   nil,
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTaskRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewTaskService(mockRepo, &worker.Worker{})

			// Act
			resp, err := svc.ListTasks(context.Background(), tt.filter)

			// Assert
			if tt.wantErr != nil {
				assert.Error(t, err)

				if errors.Is(tt.wantErr, service.ErrInvalidTaskStatus) {
					assert.True(t, errors.Is(err, service.ErrInvalidTaskStatus))
					mockRepo.AssertNotCalled(t, "List")
				} else {
					assert.EqualError(t, err, tt.wantErr.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Len(t, resp, tt.wantCount)

			for _, item := range resp {
				assert.NotEmpty(t, item.ID)
				assert.NotEmpty(t, item.Title)
				assert.NotEmpty(t, item.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func ptr[T any](t *testing.T, v T) *T {
	t.Helper()
	return &v
}
