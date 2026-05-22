package model

import (
	"github.com/google/uuid"
	"time"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text"`
	Status      TaskStatus `gorm:"type:task_status;not null;default:'todo'"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now();index:idx_tasks_created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`

	AuthorID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_tasks_author_id"`
	AssigneeID    *uuid.UUID `gorm:"type:uuid;index:idx_tasks_assignee_id"`
	BoardID       *uuid.UUID `gorm:"type:uuid;index:idx_tasks_board_id"`
	BoardColumnID *uuid.UUID `gorm:"type:uuid;index:idx_tasks_board_column_id"`
	SprintID      *uuid.UUID `gorm:"type:uuid;index:idx_tasks_sprint_id"`
	GroupID       *uuid.UUID `gorm:"type:uuid;index:idx_tasks_group_id"`

	Author      User         `gorm:"foreignKey:AuthorID"`
	Assignee    *User        `gorm:"foreignKey:AssigneeID"`
	Board       *Board       `gorm:"foreignKey:BoardID"`
	BoardColumn *BoardColumn `gorm:"foreignKey:BoardColumnID"`
	Sprint      *Sprint      `gorm:"foreignKey:SprintID"`
	Group       *Group       `gorm:"foreignKey:GroupID"`

	Watchers []User `gorm:"many2many:task_watchers;"`
}

type UserTaskStats struct {
	AsAssignee int `gorm:"column:as_assignee"`
	AsWatcher  int `gorm:"column:as_watcher"`
}
