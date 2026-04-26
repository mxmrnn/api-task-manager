package dto

import (
	"async-api-task-manager/internal/model"
	"github.com/google/uuid"
)

type TaskCreatedResponse struct {
	ID uuid.UUID `json:"id"`
}

type TaskListItemResponse struct {
	ID           uuid.UUID        `json:"id"`
	Title        string           `json:"title"`
	Status       model.TaskStatus `json:"status"`
	AssigneeID   *uuid.UUID       `json:"assignee_id,omitempty"`
	AssigneeName *string          `json:"assignee_name,omitempty"`
	BoardName    *string          `json:"board_name,omitempty"`
	ColumnName   *string          `json:"column_name,omitempty"`
}

type TaskDetailResponse struct {
	ID          uuid.UUID        `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      model.TaskStatus `json:"status"`

	Author   UserShort         `json:"author"`
	Assignee *UserShort        `json:"assignee,omitempty"`
	Board    *BoardShort       `json:"board,omitempty"`
	Column   *BoardColumnShort `json:"column,omitempty"`
	Sprint   *SprintShort      `json:"sprint,omitempty"`
	Group    *GroupShort       `json:"group,omitempty"`
}
