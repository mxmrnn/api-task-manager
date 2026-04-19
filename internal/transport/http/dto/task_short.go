package dto

import (
	"github.com/google/uuid"
)

type UserShort struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
}

type BoardShort struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BoardColumnShort struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Position int       `json:"position"`
}

type SprintShort struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GroupShort struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
