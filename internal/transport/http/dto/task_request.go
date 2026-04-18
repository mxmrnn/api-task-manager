package dto

import (
	"async-api-task-manager/internal/model"
	"github.com/google/uuid"
)

type TaskCreateRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      *model.TaskStatus `json:"status"`
	AuthorID    uuid.UUID         `json:"author_id"`
	AssigneeID  *uuid.UUID        `json:"assignee_id,omitempty"`
	BoardID     *uuid.UUID        `json:"board_id,omitempty"`
	ColumnID    *uuid.UUID        `json:"column_id,omitempty"`
	SprintID    *uuid.UUID        `json:"sprint_id,omitempty"`
	GroupID     *uuid.UUID        `json:"group_id,omitempty"`
}
