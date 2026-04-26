package dto

import (
	"async-api-task-manager/internal/model"
	"github.com/google/uuid"
)

type Filter struct {
	Author   *uuid.UUID
	Assignee *uuid.UUID
	Status   *model.TaskStatus
}
