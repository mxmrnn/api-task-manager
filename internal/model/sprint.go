package model

import (
	"github.com/google/uuid"
	"time"
)

type Sprint struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	BoardID   uuid.UUID `gorm:"type:uuid;not null;index:idx_sprints_board_id"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now();index:idx_sprints_created_at"`
}
