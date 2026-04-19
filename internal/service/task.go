package service

import (
	"async-api-task-manager/internal/storage/postgres/repository"
	"async-api-task-manager/internal/transport/http/dto"
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"async-api-task-manager/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	List(ctx context.Context, filter dto.Filter) ([]model.Task, error)
	Update(ctx context.Context, task *model.Task) error
}

type TaskService interface {
	CreateTask(ctx context.Context, req dto.TaskCreateRequest) (dto.TaskCreatedResponse, error)
	GetTaskByID(ctx context.Context, id uuid.UUID) (dto.TaskDetailResponse, error)
	ListTasks(ctx context.Context, filter dto.Filter) ([]dto.TaskListItemResponse, error)
	UpdateTask(ctx context.Context, req dto.TaskUpdateRequest, id uuid.UUID) error
}

type taskService struct {
	taskRepo TaskRepository
}

func NewTaskService(taskRepo TaskRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
	}
}

func (s *taskService) GetTaskByID(ctx context.Context, id uuid.UUID) (dto.TaskDetailResponse, error) {
	task, err := s.taskRepo.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return dto.TaskDetailResponse{}, ErrTaskNotFound
		}
		return dto.TaskDetailResponse{}, err
	}

	return toTaskDetailResponse(*task), nil
}

func toTaskDetailResponse(task model.Task) dto.TaskDetailResponse {
	var author dto.UserShort
	if task.Author.ID != uuid.Nil {
		author = dto.UserShort{
			ID:       task.Author.ID,
			FullName: task.Author.FullName,
		}
	}

	var assignee *dto.UserShort
	if task.Assignee != nil {
		assignee = &dto.UserShort{
			ID:       task.Assignee.ID,
			FullName: task.Assignee.FullName,
		}
	}

	var board *dto.BoardShort
	if task.Board != nil {
		board = &dto.BoardShort{
			ID:   task.Board.ID,
			Name: task.Board.Name,
		}
	}

	var column *dto.BoardColumnShort
	if task.BoardColumn != nil {
		column = &dto.BoardColumnShort{
			ID:       task.BoardColumn.ID,
			Name:     task.BoardColumn.Name,
			Position: task.BoardColumn.Position,
		}
	}

	var sprint *dto.SprintShort
	if task.Sprint != nil {
		sprint = &dto.SprintShort{
			ID:   task.Sprint.ID,
			Name: task.Sprint.Name,
		}
	}

	var group *dto.GroupShort
	if task.Group != nil {
		group = &dto.GroupShort{
			ID:   task.Group.ID,
			Name: task.Group.Name,
		}
	}

	return dto.TaskDetailResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,

		Author:   author,
		Assignee: assignee,
		Board:    board,
		Column:   column,
		Sprint:   sprint,
		Group:    group,
	}
}

func (s *taskService) ListTasks(ctx context.Context, filter dto.Filter) ([]dto.TaskListItemResponse, error) {
	if filter.Status != nil {
		if !isValidTaskStatus(*filter.Status) {
			return nil, ErrInvalidTaskStatus
		}
	}

	list, err := s.taskRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	res := make([]dto.TaskListItemResponse, 0, len(list))
	for _, t := range list {
		res = append(res, toTaskListItemResponse(t))
	}
	return res, nil
}

func toTaskListItemResponse(task model.Task) dto.TaskListItemResponse {
	var assigneeName *string
	if task.Assignee != nil {
		assigneeName = &task.Assignee.FullName
	}

	var boardName *string
	if task.Board != nil {
		boardName = &task.Board.Name
	}

	var columnName *string
	if task.BoardColumn != nil {
		columnName = &task.BoardColumn.Name
	}

	return dto.TaskListItemResponse{
		ID:           task.ID,
		Title:        task.Title,
		Status:       task.Status,
		AssigneeID:   task.AssigneeID,
		AssigneeName: assigneeName,
		BoardName:    boardName,
		ColumnName:   columnName,
	}
}

func (s *taskService) UpdateTask(ctx context.Context, req dto.TaskUpdateRequest, id uuid.UUID) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if utf8.RuneCountInString(title) < 3 {
			return ErrTitleTooShort
		}
		task.Title = title
	}

	if req.Description != nil {
		task.Description = strings.TrimSpace(*req.Description)
	}

	if req.Status != nil {
		status := *req.Status
		if !isValidTaskStatus(status) {
			return ErrInvalidTaskStatus
		}
		task.Status = status
	}

	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}

	if req.BoardID != nil {
		task.BoardID = req.BoardID
	}

	if req.ColumnID != nil {
		task.BoardColumnID = req.ColumnID
	}

	if req.SprintID != nil {
		task.SprintID = req.SprintID
	}

	if req.GroupID != nil {
		task.GroupID = req.GroupID
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	return nil

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
