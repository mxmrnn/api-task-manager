package repository

import (
	taskdomain "async-api-task-manager/internal/transport/http/dto"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"async-api-task-manager/internal/model"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	var task model.Task

	if err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Assignee").
		Preload("Board").
		Preload("BoardColumn").
		Preload("Sprint").
		Preload("Group").
		First(&task, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("task repository get by id: %w", ErrTaskNotFound)
		}
		return nil, fmt.Errorf("task repository get by id: %w", err)
	}

	return &task, nil
}

func (r *TaskRepository) List(ctx context.Context, filter taskdomain.Filter) ([]model.Task, error) {
	var list []model.Task

	query := r.db.WithContext(ctx).Model(&model.Task{})
	if filter.Status != nil {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Author != nil {
		query = query.Where("author_id = ?", filter.Author)
	}

	if filter.Assignee != nil {
		query = query.Where("assignee_id = ?", filter.Assignee)
	}

	if err := query.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("task repository list: %w", err)
	}

	return list, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	if err := r.db.WithContext(ctx).Save(task).Error; err != nil {
		return fmt.Errorf("task repository update: %w", err)
	}
	return nil
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return fmt.Errorf("task repository create: %w", err)
	}
	return nil
}
