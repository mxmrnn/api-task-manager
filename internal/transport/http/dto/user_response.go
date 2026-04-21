package dto

import (
	"github.com/google/uuid"
	"time"
)

type UserCreatedResponse struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
