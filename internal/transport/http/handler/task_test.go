package handler_test

import (
	"async-api-task-manager/internal/transport/http/handler"
	"context"
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
