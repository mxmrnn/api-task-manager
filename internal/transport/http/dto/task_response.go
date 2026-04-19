package dto

import (
	"async-api-task-manager/internal/model"
	"github.com/google/uuid"
)

type TaskCreatedResponse struct {
	ID uuid.UUID `json:"id"`
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
