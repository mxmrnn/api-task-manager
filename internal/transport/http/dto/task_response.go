package dto

import (
	"github.com/google/uuid"
)

type TaskCreatedResponse struct {
	ID uuid.UUID `json:"id"`
}
