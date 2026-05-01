package handler_test

import (
	"async-api-task-manager/internal/transport/http/handler"
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"async-api-task-manager/internal/service"
	"async-api-task-manager/internal/transport/http/dto"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, req dto.TaskCreateRequest) (dto.TaskCreatedResponse, error) {
	args := m.Called(ctx, req)

	var resp dto.TaskCreatedResponse
	if arg := args.Get(0); arg != nil {
		resp = arg.(dto.TaskCreatedResponse)
	}

	return resp, args.Error(1)
}

func (m *MockTaskService) GetTaskByID(ctx context.Context, id uuid.UUID) (dto.TaskDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return dto.TaskDetailResponse{}, args.Error(1)
	}
	return args.Get(0).(dto.TaskDetailResponse), args.Error(1)
}

func (m *MockTaskService) ListTasks(ctx context.Context, filter dto.Filter) ([]dto.TaskListItemResponse, error) {
	args := m.Called(ctx, filter)

	var tasks []dto.TaskListItemResponse
	if arg := args.Get(0); arg != nil {
		tasks = arg.([]dto.TaskListItemResponse)
	}

	return tasks, args.Error(1)
}

func (m *MockTaskService) UpdateTask(ctx context.Context, req dto.TaskUpdateRequest, id uuid.UUID) error {
	args := m.Called(ctx, req, id)
	return args.Error(0)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

var _ service.TaskService = (*MockTaskService)(nil)

var testRouter *gin.Engine
var mockService *MockTaskService

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	mockService = &MockTaskService{}
	testRouter = setupRouter(mockService)

	code := m.Run()
	os.Exit(code)
}

func setupRouter(svc *MockTaskService) *gin.Engine {
	r := gin.New()
	h := handler.NewTaskHandler(svc)

	r.GET("/health", h.Health)
	r.GET("/tasks/:id", h.GetTaskByID)
	r.GET("/tasks", h.ListTasks)
	r.POST("/tasks", h.CreateTask)
	r.PUT("/tasks/:id", h.UpdateTask)
	r.DELETE("/tasks/:id", h.DeleteTask)

	return r
}

func TestTaskHandler_ListTasks(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		mockSetup  func()
		wantStatus int
		wantBody   string
	}{
		{
			name:  "success list tasks",
			query: "?status=todo&author_id=123e4567-e89b-12d3-a456-426614174000",
			mockSetup: func() {
				mockService.On("ListTasks", mock.Anything, mock.AnythingOfType("dto.Filter")).
					Return([]dto.TaskListItemResponse{
						{ID: uuid.New(), Title: "Task 1", Status: "todo"},
						{ID: uuid.New(), Title: "Task 2", Status: "todo"},
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   "Task 1",
		},
		{
			name:       "non valid author_id",
			query:      "?author_id=non_valid-uuid",
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid author_id",
		},
		{
			name:       "non valid assignee_id",
			query:      "?assignee_id=123",
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid assignee_id",
		},
		{
			name:  "err service",
			query: "?status=todo",
			mockSetup: func() {
				mockService.On("ListTasks", mock.Anything, mock.Anything).Return([]dto.TaskListItemResponse{}, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed to list tasks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			t.Cleanup(func() {
				mockService.AssertExpectations(t)
			})

			req := httptest.NewRequest("GET", "/tasks"+tt.query, nil)
			w := httptest.NewRecorder()

			testRouter.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantBody != "" {
				assert.Contains(t, w.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestTaskHandler_CreateTask(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		mockSetup    func()
		wantStatus   int
		wantContains string
	}{
		{
			name: "success create task",
			body: `{
				"title": "task 1",
				"description": "description 1",
				"authorId": "123e4567-e89b-12d3-a456-426614174000"
			}`,
			mockSetup: func() {
				mockService.On("CreateTask", mock.Anything, mock.Anything).
					Return(dto.TaskCreatedResponse{ID: uuid.New()}, nil)
			},
			wantStatus:   http.StatusCreated,
			wantContains: `"id":`,
		},
		{
			name: "err title to short",
			body: `{"title": "ab"}`,
			mockSetup: func() {
				mockService.On("CreateTask", mock.Anything, mock.Anything).
					Return(dto.TaskCreatedResponse{}, service.ErrTitleTooShort)
			},
			wantStatus:   http.StatusBadRequest,
			wantContains: "title must be at least 3 characters",
		},
		{
			name: "err status not valid",
			body: `{"title": "title one", "status": "not valid"}`,
			mockSetup: func() {
				mockService.On("CreateTask", mock.Anything, mock.Anything).
					Return(dto.TaskCreatedResponse{}, service.ErrInvalidTaskStatus)
			},
			wantStatus:   http.StatusBadRequest,
			wantContains: "invalid task status",
		},
		{
			name:         "err non correct JSON",
			body:         `{"title": "title one", "status": non valid}`,
			mockSetup:    func() {},
			wantStatus:   http.StatusBadRequest,
			wantContains: "invalid request body",
		},
		{
			name: "err from service(внутренняя)",
			body: `{"title": "title one"}`,
			mockSetup: func() {
				mockService.On("CreateTask", mock.Anything, mock.Anything).
					Return(dto.TaskCreatedResponse{}, assert.AnError)
			},
			wantStatus:   http.StatusInternalServerError,
			wantContains: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil

			// Arrange
			tt.mockSetup()

			t.Cleanup(func() {
				mockService.AssertExpectations(t)
			})

			req := httptest.NewRequest("POST", "/tasks", bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			// Act
			testRouter.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantContains != "" {
				assert.Contains(t, w.Body.String(), tt.wantContains)
			}

		})
	}
}
