package service

import (
	"async-api-task-manager/internal/transport/http/dto"
	"context"
	"strings"
	"unicode/utf8"

	"async-api-task-manager/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
}

type TaskService interface {
	CreateTask(ctx context.Context, req dto.TaskCreateRequest) (dto.TaskCreatedResponse, error)
}

type taskService struct {
	taskRepo TaskRepository
}

func NewTaskService(taskRepo TaskRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
	}
}

func (s *taskService) CreateTask(ctx context.Context, req dto.TaskCreateRequest) (dto.TaskCreatedResponse, error) {
	title := strings.TrimSpace(req.Title)
	if utf8.RuneCountInString(title) < 3 {
		return dto.TaskCreatedResponse{}, ErrTitleTooShort
	}
	status := model.TaskStatusTodo

	if req.Status != nil {
		status = *req.Status
	}

	if !isValidTaskStatus(status) {
		return dto.TaskCreatedResponse{}, ErrInvalidTaskStatus
	}

	task := model.Task{
		Title:         title,
		Description:   strings.TrimSpace(req.Description),
		Status:        status,
		AuthorID:      req.AuthorID,
		AssigneeID:    req.AssigneeID,
		BoardID:       req.BoardID,
		BoardColumnID: req.ColumnID,
		SprintID:      req.SprintID,
		GroupID:       req.GroupID,
	}

	if err := s.taskRepo.Create(ctx, &task); err != nil {
		return dto.TaskCreatedResponse{}, err
	}

	return toTaskCreatedResponse(task), nil
}

func isValidTaskStatus(status model.TaskStatus) bool {
	switch status {
	case model.TaskStatusTodo, model.TaskStatusInProgress, model.TaskStatusDone:
		return true
	default:
		return false
	}
}

func toTaskCreatedResponse(task model.Task) dto.TaskCreatedResponse {
	return dto.TaskCreatedResponse{
		ID: task.ID,
	}
}
